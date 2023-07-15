package utils

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Email    string
	Hash     string
	TokenKey string
}

type EnvTest struct {
	Email    string
	Password string
	Hash     string
}

func EnvVariables() (Env, error) {
	var env Env
	err := godotenv.Load("../.env")
	if err != nil {
		return env, err
	}
	env.Email = os.Getenv("EMAIL")
	env.Hash = os.Getenv("PW_HASH")
	env.TokenKey = os.Getenv("TOKEN_KEY")
	return env, nil
}

func EnvTestVariables() (EnvTest, error) {
	var envTest EnvTest
	err := godotenv.Load("../.env.dev")
	if err != nil {
		return envTest, err
	}
	envTest.Email = os.Getenv("TEST_EMAIL")
	envTest.Password = os.Getenv("TEST_PW")
	envTest.Hash = os.Getenv("TEST_HASH")
	return envTest, nil
}
