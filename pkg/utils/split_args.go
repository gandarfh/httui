package utils

import (
	"regexp"
	"strings"
)

// Will match: word.word='some value to this'
// Will match: word="some value to this"
// Will match: word=42
// Implementation: https://regexr.com/6to57
var re = regexp.MustCompile(`(((\w\.?\w?)*?)=(("|')(.*?)("|')|\w*))`)

func SplitArgs(args string) []string {
	command := strings.Trim(re.ReplaceAllString(args, ""), " ")

	commands := []string{command}
	keysAndValues := []string{}

	if IsKeyValue(args) {
		for _, item := range re.FindAllString(args, -1) {
			newValue := strings.Replace(item, "'", "", -1)
			newValue = strings.Replace(newValue, `"`, "", -1)

			keysAndValues = append(keysAndValues, newValue)
		}

	}

	commands = append(commands, keysAndValues...)
	return commands
}

func IsKeyValue(command string) bool {
	has := re.MatchString(command)

	return has
}
