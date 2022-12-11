package utils

import (
	"regexp"
	"strings"

	"github.com/gandarfh/httui/internal/repositories"
)

var re_env = regexp.MustCompile(`\$+.\S+`)

func ReadEnv(key string) (repositories.Env, error) {
	repo, _ := repositories.NewEnvs()

	key = strings.Replace(key, "$", "", 1)
	env, err := repo.FindByKey(key)

	return env, err
}

func HaveEnv(raw string) bool {
	return re_env.Match([]byte(raw))
}

func ReplaceByEnv(raw string) string {
	listOfEnvs := re_env.FindAllString(raw, -1)

	for _, item := range listOfEnvs {
		if env, err := ReadEnv(item); err != nil {
			raw = strings.ReplaceAll(raw, item, item)
		} else {
			raw = strings.ReplaceAll(raw, item, env.Value)
		}
	}

	return raw
}
