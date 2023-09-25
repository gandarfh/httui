package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gandarfh/httui/external/database"
	"github.com/gandarfh/httui/internal/repositories"
	"gopkg.in/yaml.v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func processYMLForRequests() {
	db := database.Client

	folderPath := "../insomnia-to-httui/insomnia/Request"
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		panic("Erro ao ler a pasta: " + err.Error())
	}

	// Criação inicial dos grupos
	parents := make(map[string][]repositories.Request)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yml" {
			filePath := filepath.Join(folderPath, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println("Erro ao ler o arquivo", filePath, ":", err)
				continue
			}

			var requestData map[string]interface{}
			err = yaml.Unmarshal(data, &requestData)
			if err != nil {
				fmt.Println("Erro ao analisar o arquivo", filePath, ":", err)
				continue
			}

			request := createRequest(db, requestData)
			if request != nil {
				parents[requestData["parentId"].(string)] = append(parents[requestData["parentId"].(string)], *request)
			}
		}
	}

	updateRequestsParents(db, parents)
}

func convertToParams(value interface{}) datatypes.JSONType[[]map[string]string] {
	if data := value.([]interface{}); len(data) == 0 {
		return datatypes.NewJSONType([]map[string]string{})
	}

	data := []map[string]string{}

	for _, key := range value.([]interface{}) {
		param := key.(map[interface{}]interface{})

		data = append(data, map[string]string{
			param["name"].(string): param["value"].(string),
		})

	}

	return datatypes.NewJSONType(data)
}

func convertToBody(value interface{}) datatypes.JSONType[map[string]interface{}] {
	data := value.(map[interface{}]interface{})["text"]

	if data == nil || value == nil {
		return datatypes.NewJSONType(map[string]interface{}{})
	}

	body := map[string]interface{}{}
	json.Unmarshal([]byte(data.(string)), &body)

	return datatypes.NewJSONType(body)
}

func createRequest(db *gorm.DB, data map[string]interface{}) *repositories.Request {
	fmt.Println(data["body"], "\n", convertToBody(data["body"]))

	request := &repositories.Request{
		Type:             "request",
		Name:             data["name"].(string),
		ExternalId:       data["_id"].(string),
		ExternalParentId: data["parentId"].(string),
		Method:           data["method"].(string),
		Endpoint:         data["url"].(string),
		QueryParams:      convertToParams(data["parameters"]),
		Headers:          convertToParams(data["headers"]),
		Body:             convertToBody(data["body"]),
	}

	if err := db.Create(request).Error; err != nil {
		fmt.Println("Erro ao inserir no banco de dados:", err)
		return nil
	}

	return request
}

func updateRequestsParents(db *gorm.DB, parents map[string][]repositories.Request) {
	for name, requests := range parents {
		parent := repositories.Request{}
		if err := db.Model(&parent).Where("type = ?", "group").Where("external_id = ?", name).First(&parent).Error; err != nil {
			fmt.Println("Erro ao puxar parent request", err)
			continue
		}

		for _, request := range requests {
			if parent.ID != request.ID {
				request.ParentID = &parent.ID
				if err := db.Save(&request).Error; err != nil {
					fmt.Println("Erro ao atualizar o grupo:", err)
				}
			}
		}
	}
}
