package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gandarfh/httui/external/database"
	"github.com/gandarfh/httui/internal/repositories"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

func processYMLForGroups() {
	db := database.Client

	folderPath := "../insomnia-to-httui/insomnia/RequestGroup"
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

			var requestGroupData map[string]interface{}
			err = yaml.Unmarshal(data, &requestGroupData)
			if err != nil {
				fmt.Println("Erro ao analisar o arquivo", filePath, ":", err)
				continue
			}

			groupRequest := createGroupRequest(db, requestGroupData)
			if groupRequest != nil {
				parents[requestGroupData["parentId"].(string)] = append(parents[requestGroupData["parentId"].(string)], *groupRequest)
			}
		}
	}

	updateGroupParents(db, parents)

}

func createGroupRequest(db *gorm.DB, data map[string]interface{}) *repositories.Request {
	groupRequest := &repositories.Request{
		Type:             "group",
		Name:             data["name"].(string),
		ExternalId:       data["_id"].(string),
		ExternalParentId: data["parentId"].(string),
	}

	if err := db.Create(groupRequest).Error; err != nil {
		fmt.Println("Erro ao inserir no banco de dados:", err)
		return nil
	}

	return groupRequest
}

func updateGroupParents(db *gorm.DB, parents map[string][]repositories.Request) {
	for name, groups := range parents {
		parent := repositories.Request{}
		if err := db.Model(&parent).Where("external_id = ?", name).First(&parent).Error; err != nil {
			fmt.Println("Erro ao puxar parent group", err)
			continue
		}

		for _, group := range groups {
			if parent.ID != group.ID {
				group.ParentID = &parent.ID
				if err := db.Save(&group).Error; err != nil {
					fmt.Println("Erro ao atualizar o grupo:", err)
				}
			}
		}
	}
}
