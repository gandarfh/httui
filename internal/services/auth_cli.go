package services

import (
	"encoding/json"
	"io"

	"github.com/gandarfh/httui/pkg/client"
)

func AuthCLI(name string) string {
	body := map[string]string{
		"Name": name,
	}
	data, _ := json.Marshal(body)
	api, _ := client.Post("http://localhost:6000/auth/cli").Body(data).Exec()

	var response map[string]interface{}
	readbody, _ := io.ReadAll(api.Body)
	json.Unmarshal(readbody, &response)

	return response["url"].(string)
}
