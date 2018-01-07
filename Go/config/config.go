package config

import (
	"errors"
	"github.com/joho/godotenv"
)

func FunInitConfig() error {
	err := godotenv.Load("environment.env")
	if err != nil {
		err = errors.New("FunInitConfig: couldn't initialize environment -> " + err.Error())
		return err
	}
	return nil
}
