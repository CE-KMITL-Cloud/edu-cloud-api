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
	REALM       = "@pve"
	AUTH_COOKIE = "PVEAuthCookie"
	CSRF_TOKEN  = "CSRFPreventionToken"
	URL_ENCODED = "application/x-www-form-urlencoded"
	Gigabyte    = 1073741824 // Gigabyte : 1024^3
	Megabyte    = 1048576    // Megabyte : 1024^2
	Byte        = 1024       // Byte : 1024^1
	MatchNumber = `:(\d+)`
	WorkerNode  = `work-[-]?\d[\d,]*[\.]?[\d{2}]*`
	ENV_PATH    = ".env"

	// Create VM's Configuration
	SCSIHW = "virtio-scsi-pci"
	NET0   = "virtio,bridge=vmbr0,firewall=1"
	ISO    = "cephfs:iso/"
	SOCKET = 1
	ONBOOT = 1

	// DBs
	ADMIN   = "admin"
	STUDENT = "student"
	FACULTY = "faculty"
)

// GBtoByte - Converter from GB to Byte
func GBtoByte(input uint64) uint64 {
	return input * Gigabyte
}

// GBtoByteFloat - Converter from GB to Byte Float
func GBtoByteFloat(input float64) uint64 {
	return uint64(input * Gigabyte)
}

// BytetoGB - Converter from Byte to GB
func BytetoGB(input uint64) float64 {
	return float64(input) / Gigabyte
}

// MBtoByte - Converter from MB to Byte
func MBtoByte(input uint64) uint64 {
	return input * Megabyte
}

// MBtoGB - Converter from MB to GB
func MBtoGB(input uint64) float64 {
	return float64(input / Byte)
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
	godotenv.Load(ENV_PATH)
	return os.Getenv(item)
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

// SendRequestUsingToken - Constructing HTTP client and Sending request using token
func SendRequestUsingToken(httpMethod, url string, data url.Values) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", GetFromENV("PROXMOX_API_KEY"))
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

// FilterList - compare two string list
func FilterList(list1, list2 []string) []string {
	filtered := make([]string, 0, len(list1))
	for _, item := range list1 {
		if !Contains(list2, item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
