package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ArgsFormat(args []string) (map[string]any, error) {
	mapArgs := map[string]any{}

	for _, arg := range args {
		item := strings.Split(arg, "=")

		if len(item) == 2 {
			key, value := item[0], item[1]

			// convert string to int
			newValue, err := strconv.Atoi(value)

			if err == nil {
				// convert to int value
				mapArgs[key] = newValue
				continue
			}

			mapArgs[key] = value
			continue

		}

		return nil, fmt.Errorf("Key and value not expected. You provide: %s.\nTry something like: key=value", arg)
	}

	return mapArgs, nil
}

func IsInt(value string) bool {
	matched, _ := regexp.Match("/[0-9]/g", []byte(value))

	return matched
}
