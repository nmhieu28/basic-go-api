package environment

import (
	"os"
)

type Environment string

var (
	Development Environment = "development"
	Production  Environment = "production"
	Test        Environment = "test"
)

func GetEnvironment() Environment {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = string(Development)
	}

	return Environment(env)
}

func (env Environment) IsDevelopment() bool {
	return env == Development
}

func (env Environment) IsProduction() bool {
	return env == Production
}

func (env Environment) IsTest() bool {
	return env == Test
}

func (env Environment) String() string {
	return string(env)
}
