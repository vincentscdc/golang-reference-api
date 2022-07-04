package configuration

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type MissingEnvConfigError struct {
	env string
	err error
}

func (mece MissingEnvConfigError) Error() string {
	return fmt.Sprintf("missing config %s: %v", mece.env, mece.err)
}

func (mece MissingEnvConfigError) Is(target error) bool {
	// nolint: errorlint, goerr113 // https://github.com/polyfloyd/go-errorlint/issues/11#issuecomment-873072046
	return target == MissingEnvConfigError{}
}

type MissingBaseConfigError struct {
	err error
}

func (mbce MissingBaseConfigError) Error() string {
	return fmt.Sprintf("missing base config: %v", mbce.err)
}

func (mbce MissingBaseConfigError) Is(target error) bool {
	// nolint: errorlint, goerr113 // https://github.com/polyfloyd/go-errorlint/issues/11#issuecomment-873072046
	return target == MissingBaseConfigError{}
}

func GetConfig(configPath, currEnv string) (Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath+"/base.yaml", &cfg); err != nil {
		return cfg, MissingBaseConfigError{err: err}
	}

	if err := cleanenv.ReadConfig(fmt.Sprintf(configPath+"/%s.yaml", currEnv), &cfg); err != nil {
		return cfg, MissingEnvConfigError{env: currEnv, err: err}
	}

	return cfg, nil
}
