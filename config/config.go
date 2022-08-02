package config

import (
	"fmt"
	"github/gandarfh/httui/utils"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type connection = interface{}

func Connect[T interface{}]() T {
	file, err := os.ReadFile("config.yml")
	var data T

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error when try read config.yml file:\n")

		return data
	}

	err = yaml.Unmarshal(file, &data)

	if err != nil {
		utils.MsgError(err, "Bad syntax in config.yml. \n", "Syntax must be like: `Config struct`")
		return data
	}

	return data
}
