package utils

import (
	"strings"

	"github.com/gandarfh/maid-san/internal/commands/envs/repository"
	"github.com/gandarfh/maid-san/pkg/errors"
)

func ReadEnv(key string) (*repository.Envs, error) {
	repo, err := repository.NewEnvsRepo()

	env := &repository.Envs{
		Key:   key,
		Value: key,
	}

	if err != nil {
		return nil, errors.InternalServer("Func IsEnvKey: error when try connect to database!")
	}

	// Valid if have $prefix
	if !strings.Contains(key, "$") {
		return env, nil
	}

	key = strings.Replace(key, "$", "", 1)

	env, err = repo.FindByKey(key)

	if err != nil {
		return &repository.Envs{
			Key:   key,
			Value: key,
		}, nil
	}

	return env, nil
}
