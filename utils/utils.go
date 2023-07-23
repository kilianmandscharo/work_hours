package utils

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
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

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	parentDir := filepath.Dir(currentDir)
	path := filepath.Join(parentDir, ".env")

	err := godotenv.Load(path)
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

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	parentDir := filepath.Dir(currentDir)
	path := filepath.Join(parentDir, ".env.dev")

	err := godotenv.Load(path)
	if err != nil {
		return envTest, err
	}
	envTest.Email = os.Getenv("TEST_EMAIL")
	envTest.Password = os.Getenv("TEST_PW")
	envTest.Hash = os.Getenv("TEST_HASH")
	return envTest, nil
}
