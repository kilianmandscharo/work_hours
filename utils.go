package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	email    string
	hash     string
	tokenKey string
}

type EnvTest struct {
	email    string
	password string
	hash     string
}

func envVariables() (Env, error) {
	var env Env
	err := godotenv.Load(".env")
	if err != nil {
		return env, err
	}
	env.email = os.Getenv("EMAIL")
	env.hash = os.Getenv("PW_HASH")
	env.tokenKey = os.Getenv("TOKEN_KEY")
	return env, nil
}

func envTestVariables() (EnvTest, error) {
	var envTest EnvTest
	err := godotenv.Load(".env.dev")
	if err != nil {
		return envTest, err
	}
	envTest.email = os.Getenv("TEST_EMAIL")
	envTest.password = os.Getenv("TEST_PW")
	envTest.hash = os.Getenv("TEST_HASH")
	return envTest, nil
}
