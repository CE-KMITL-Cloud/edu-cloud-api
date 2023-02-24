// Package config - for utils function
package config

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

const (
	AUTH_COOKIE = "PVEAuthCookie"
	CSRF_TOKEN  = "CSRFPreventionToken"
	URL_ENCODED = "application/x-www-form-urlencoded"
	Gigabyte    = 1073741824 // Gigabyte : 1024^3
	Megabyte    = 1048576    // Megabyte : 1024^2
	MatchNumber = `:(\d+)`
	WorkerNode  = `work-[-]?\d[\d,]*[\.]?[\d{2}]*`
)

// GBtoByte - Converter from GB to Byte
func GBtoByte(input uint64) uint64 {
	return input * Gigabyte
}

// MBtoByte - Converter from MB to Byte
func MBtoByte(input uint64) uint64 {
	return input * Megabyte
}

// GreaterOrEqual - Compare sets of VM spec {cpu, mem (byte), disk (byte)}
func GreaterOrEqual(cpuA, cpuB float64, memA, memB uint64, diskA, diskB uint64) bool {
	if cpuA == cpuB && memA == memB && diskA == diskB {
		return false
	} else if cpuA > cpuB && memA >= memB && diskA >= diskB {
		return true
	} else if cpuA >= cpuB && memA > memB && diskA >= diskB {
		return true
	} else if cpuA >= cpuB && memA >= memB && diskA > diskB {
		return true
	} else {
		return false
	}
}

// GetFromENV - get item from .env
func GetFromENV(item string) string {
	godotenv.Load(".env")
	return os.Getenv(item)
}

// GetListFromENV - get item list from .env
func GetListFromENV(item []string) []string {
	godotenv.Load(".env")
	var list []string
	for i := 0; i < len(item); i++ {
		value := os.Getenv(item[i])
		list = append(list, value)
	}
	return list
}

// GetURL - Constructing Proxmox's API URL
func GetURL(query string) string {
	// hostURL := GetFromENV("PROXMOX_HOST")
	hostURL := os.Getenv("PROXMOX_HOST")
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = query
	return u.String()
}

// SendRequest - Constructing HTTP client and Sending request
func SendRequest(httpMethod, url string, data url.Values, cookies model.Cookies) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	if data != nil {
		req.Header.Add("Content-Type", URL_ENCODED)
	}
	req.AddCookie(&cookies.Cookie)
	req.Header.Add(CSRF_TOKEN, cookies.CSRFPreventionToken.Value)

	client := &http.Client{}
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return nil, sendErr
	}
	return resp, nil
}

// SendRequestWithErr - Constructing HTTP client and Sending request
func SendRequestWithErr(httpMethod, url string, data url.Values, cookies model.Cookies) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	if data != nil {
		req.Header.Add("Content-Type", URL_ENCODED)
	}
	req.AddCookie(&cookies.Cookie)
	req.Header.Add(CSRF_TOKEN, cookies.CSRFPreventionToken.Value)

	client := &http.Client{}
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return nil, sendErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error: with status", resp.Status)
		return nil, errors.New(resp.Status)
	}
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	// log.Println(string(respBody))
	return respBody, nil
}

// SendRequestWithoutCookie - Constructing HTTP client and Sending request without cookies
func SendRequestWithoutCookie(httpMethod, url string, data url.Values) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	if data != nil {
		req.Header.Add("Content-Type", URL_ENCODED)
	}

	client := &http.Client{}
	resp, sendErr := client.Do(req)
	if sendErr != nil {
		return nil, sendErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error: with status", resp.Status)
		return nil, errors.New(resp.Status)
	}
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	// log.Println(string(respBody))
	return respBody, nil
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
