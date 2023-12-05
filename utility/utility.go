package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/jwalton/go-supportscolor"
)

type tsHttp struct {
	Protocol string `json:"protocol"`
	Host 	 string `json:"host"`
	HttpPort string `json:"http_port"`
}

type tsDatabase struct {
	Hostname     string `json:"hostname"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
}

type tsStorage struct {
	BucketName  string `json:"bucket_name"`
	Credentials string `json:"credentials"`
}

type Configuration struct {
	Http     tsHttp     `json:"http"`
	Database tsDatabase `json:"database"`
	AppPath	 string     `json:"app_path"`
	Storage  tsStorage  `json:"storage"`
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func LoadApplicationConfiguration(removeSuffixPath string) (config Configuration, err error) {
	defer RecoverError()

	config.AppPath, err = os.Getwd()
	if err != nil {
		return
	}

	if removeSuffixPath != "" {
		config.AppPath = strings.TrimSuffix(config.AppPath, removeSuffixPath)
	}

	configPath := filepath.Join(config.AppPath, "config.json")
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	if err != nil {
		return
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &config)
		return
	}
}

func RecoverError() {
	if r := recover(); r != nil {
		PrintConsole(fmt.Sprintf("[ERROR][RECOVER]=> %v", r), "error")
	}
}

func PrintConsole(strPrint string, strStatus string) {
	defer RecoverError()

	if supportscolor.Stdout().SupportsColor {
		if strings.ToLower(strings.TrimSpace(strStatus)) == "info" {
			fmt.Println(Green + "[INFO]> " + strPrint + Reset)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "error" {
			fmt.Println(Red + "[ERROR]> " + strPrint + Reset)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "warning" {
			fmt.Println(Yellow + "[WARNING]> " + strPrint + Reset)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "logo" {
			fmt.Println(Red + strPrint + Reset)
		} else {
			fmt.Println(Reset + strPrint + Reset)
		}
	} else {
		if strings.ToLower(strings.TrimSpace(strStatus)) == "info" {
			fmt.Println("[INFO]> " + strPrint)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "error" {
			fmt.Println("[ERROR]> " + strPrint)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "warning" {
			fmt.Println("[WARNING]> " + strPrint)
		} else if strings.ToLower(strings.TrimSpace(strStatus)) == "logo" {
			fmt.Println(strPrint)
		} else {
			fmt.Println(strPrint)
		}
	}
}