package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gandarfh/httui/internal/repositories"
)

var re_env = regexp.MustCompile(`{{ _.\w* }}`)
var re_response = regexp.MustCompile(`{% response '\w+(\.\w+)*' '\d*' %}`)
var re_response_fields = regexp.MustCompile(`'\w+(\.\w+)*'`)

// {% response 'field' '1' %}

func ReadEnv(key string, workspaceId uint) (repositories.Env, error) {
	repo := repositories.NewEnvs()

	key = strings.Replace(key, "{{ _.", "", 1)
	key = strings.Replace(key, " }}", "", 1)

	env, err := repo.FindByKey(key, workspaceId)

	return env, err
}

func getField(payload map[string]interface{}, fields []string) (interface{}, error) {
	if len(fields) == 0 {
		return payload, nil
	}

	field := fields[0]
	value, found := payload[field]
	if !found {
		return nil, errors.New("Not found")
	}

	if submap, ok := value.(map[string]interface{}); ok {
		return getField(submap, fields[1:])
	}

	return value, nil
}

func ReadResponse(key string) (string, error) {
	infos := re_response_fields.FindAllString(key, -1)

	infos[0] = strings.Replace(infos[0], "'", "", -1)
	reqId := strings.Replace(infos[1], "'", "", -1)

	fields := strings.Split(infos[0], ".")

	repo := repositories.NewResponse()
	response, err := repo.FindOne(reqId)

	data := map[string]interface{}{}
	Convert(response, &data)

	value, _ := getField(data, fields)

	return fmt.Sprint(value), err
}

type TransformFunc func(text string) string

func ReplaceByOperator(raw string, workspaceId uint, transforms ...TransformFunc) string {
	listOfEnvs := re_env.FindAllString(raw, -1)
	listOfResponses := re_response.FindAllString(raw, -1)

	for _, item := range listOfEnvs {
		if env, err := ReadEnv(item, workspaceId); err != nil {
			raw = strings.ReplaceAll(raw, item, item)
		} else {
			raw = strings.ReplaceAll(raw, item, env.Value)
		}
	}

	for _, item := range listOfResponses {
		if data, err := ReadResponse(item); err != nil {
			raw = strings.ReplaceAll(raw, item, item)
		} else {
			raw = strings.ReplaceAll(raw, item, data)
		}
	}

	return raw
}
