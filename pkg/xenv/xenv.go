package xenv

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
)

func Loads[T any](t *T) (*T, error) {
	if t == nil {
		t = new(T)
	}

	if err := env.Parse(t); err != nil {
		return nil, fmt.Errorf("xenv: failed to parse environment variables: %w", err)
	}

	return t, nil
}
