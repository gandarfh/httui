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

func SplitArgs(command string) []string {
	commands := strings.Split(command, " ")

	if hasQuotes(command) {
		removedQuotes := []string{}

		for _, item := range re.FindAllString(command, -1) {
			newValue := strings.Replace(item, "'", "", -1)
			newValue = strings.Replace(newValue, `"`, "", -1)

			removedQuotes = append(removedQuotes, newValue)

		}

		commands = append(commands[:1], removedQuotes...)
		return commands

	}

	return commands
}

func hasQuotes(command string) bool {
	has := re.MatchString(command)

	return has
}
