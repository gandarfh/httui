package convert

import "encoding/json"

func MapToStruct(value map[string]any, decode any) error {
	jsonBody, err := json.Marshal(value)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonBody, decode); err != nil {
		return err
	}

	return nil
}
