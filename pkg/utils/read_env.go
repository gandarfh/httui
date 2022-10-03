package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/envs/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
)

var re_env = regexp.MustCompile(`\$+.\S+`)

func ReadEnv(key string) (*repository.Envs, error) {
	repo, err := repository.NewEnvsRepo()

	env := &repository.Envs{
		Key:   key,
		Value: key,
	}

	if err != nil {
		return nil, errors.InternalServer("Func IsEnvKey: error when try connect to database!")
	}

	key = strings.Replace(key, "$", "", 1)

	env, err = repo.FindByKey(key)

	if err != nil {
		return nil, errors.NotFoundError()
	}

	return env, nil
}

func HaveEnv(raw string) bool {
	return re_env.Match([]byte(raw))
}

func ReplaceByEnv(raw string) (string, error) {
	listOfEnvs := re_env.FindAllString(raw, -1)

	for _, item := range listOfEnvs {
		if env, err := ReadEnv(item); err != nil {

			key := strings.Replace(item, "$", "", 1)
			fmt.Printf("[warn] - Not found env [%s]. Create it with the following command:\n", item)
			fmt.Printf(`envs create key="%s" value="some value here"%s`, key, "\n\n")

			raw = strings.ReplaceAll(raw, item, item)
		} else {
			raw = strings.ReplaceAll(raw, item, env.Value)
		}
	}

	return raw, nil
}
