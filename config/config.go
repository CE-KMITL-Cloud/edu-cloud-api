// Package config - for utils function
package config

import (
	"log"
	"net/http"

	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

const (
	AUTH_COOKIE = "PVEAuthCookie"
	CSRF_TOKEN  = "CSRFPreventionToken"
	URL_ENCODED = "application/x-www-form-urlencoded"
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

// GetCookies - Getting PVE cookie & CSRF Prevention Token
func GetCookies(c *fiber.Ctx) model.Cookies {
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  AUTH_COOKIE,
			Value: c.Cookies(AUTH_COOKIE),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  CSRF_TOKEN,
			Value: c.Cookies(CSRF_TOKEN),
		},
	}
	return cookies
}

// Contains - check string in list
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
