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

func processYMLForEnvs() {
	db := database.Client

	folderPath := "../insomnia-to-httui/insomnia/Environment"
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		panic("Erro ao ler a pasta: " + err.Error())
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yml" {
			filePath := filepath.Join(folderPath, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println("Erro ao ler o arquivo", filePath, ":", err)
				continue
			}

			var workspaceData map[string]interface{}
			err = yaml.Unmarshal(data, &workspaceData)
			if err != nil {
				fmt.Println("Erro ao analisar o arquivo", filePath, ":", err)
				continue
			}

			workspace := createWorkspace(db, workspaceData)
			envData := workspaceData["data"].(map[interface{}]interface{})

			for key, value := range envData {
				data := map[string]string{
					"key":   fmt.Sprint(key),
					"value": fmt.Sprint(value),
				}

				createEnv(db, workspace.ID, data)

			}

		}
	}
}

func createEnv(db *gorm.DB, workspaceId uint, data map[string]string) *repositories.Env {
	groupRequest := &repositories.Env{
		WorkspaceId: workspaceId,
		Key:         data["key"],
		Value:       data["value"],
	}

	if err := db.Create(groupRequest).Error; err != nil {
		fmt.Println("Erro ao inserir no banco de dados:", err)
		return nil
	}

	return groupRequest
}

func createWorkspace(db *gorm.DB, data map[string]interface{}) *repositories.Workspace {
	workspace := &repositories.Workspace{
		Name: data["name"].(string),
	}

	if err := db.Create(workspace).Error; err != nil {
		fmt.Println("Erro ao inserir no banco de dados:", err)
		return nil
	}

	return workspace
}
