// Package config - for utils function
package config

import (
	"log"

	"github.com/spf13/viper"
)

// getFromENV - get item from .env
func getFromENV(item string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(item).(string)
	if !ok {
		log.Fatalf("Error while getting item : %s", value)
	}
	return value
}

// Collecting variables from .env
// AUTHNAME := getFromENV("AUTHNAME")
