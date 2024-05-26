package debug

import (
	"encoding/json"
	"log"
)

func Debug(a ...interface{}) {
	for _, item := range a {
		value, _ := json.MarshalIndent(item, "", "  ")
		log.Println(string(value))
	}
}
