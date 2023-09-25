package main

import (
	"fmt"

	"github.com/gandarfh/httui/external/database"
	"github.com/gandarfh/httui/internal/repositories"
)

func main() {
	if err := database.SqliteConnection(); err != nil {
		panic("Erro ao abrir o banco de dados: " + err.Error())
	}

	db := database.Client

	db.AutoMigrate(&repositories.Workspace{})
	db.AutoMigrate(&repositories.Env{})
	db.AutoMigrate(&repositories.Request{})

	processYMLForGroups()
	processYMLForRequests()
	// processYMLForEnvs()

	fmt.Println("Dados populados no banco de dados SQLite com sucesso!")
}
