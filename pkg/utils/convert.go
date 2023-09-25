package utils

import "encoding/json"

func Convert(source interface{}, destination interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, destination)
}
