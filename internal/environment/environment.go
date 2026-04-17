package environment

import "fmt"

type Environment string

const (
	DevelopmentEnv Environment = "development"
	TestingEnv     Environment = "testing"
	ProductionEnv  Environment = "production"
)

var currentEnv = DevelopmentEnv

func MustInit(env string) {
	Init(MustFromString(env))
}

func Init(env Environment) {
	currentEnv = env
}

func Get() Environment {
	return currentEnv
}

func MustFromString(envStr string) Environment {
	env, err := FromString(envStr)
	if err != nil {
		panic(err)
	}
	return env
}

func FromString(envStr string) (Environment, error) {
	switch envStr {
	case string(DevelopmentEnv):
		return DevelopmentEnv, nil
	case string(TestingEnv):
		return TestingEnv, nil
	case string(ProductionEnv):
		return ProductionEnv, nil
	}
	return "", fmt.Errorf("Unknown environment %s", envStr)
}
