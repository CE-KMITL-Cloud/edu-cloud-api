// Package config - for utils function
package config

import (
	"log"

	"github.com/spf13/viper"
)

// GetFromENV - get item from .env
func GetFromENV(item string) string {
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

// GetListFromENV - get item list from .env
func GetListFromENV(item []string) []string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	var list []string
	for i := 0; i < len(item); i++ {
		value, ok := viper.Get(item[i]).(string)
		if !ok {
			log.Fatalf("Error while getting item : %s", value)
		}
		list = append(list, value)
	}
	return list
}
