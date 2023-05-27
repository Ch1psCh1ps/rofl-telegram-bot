package cmd

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() {
	realPath, errGetWd := os.Getwd()

	if errGetWd != nil {
		panic(errGetWd)
	}

	if err := godotenv.Load(realPath + "/.env"); err != nil {
		panic(err)
	}
}
