package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Configuration struct {
	IgnoredPaths []string `json:"ignore"`
}

func Init() *Configuration {
	file := checkFileExistAndCreate()
	defer file.Close()

	cfg := &Configuration{}
	err := json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (cfg *Configuration) IsIgnored(entryPath string) bool {
	entryPath = entryPath[2:]
	for _, path := range cfg.IgnoredPaths {
		if strings.HasPrefix(entryPath, path) {
			return true
		}
	}
	return false
}

func checkFileExistAndCreate() *os.File {
	_, err := os.Stat(".bigpicture.json")
	if err != nil {
		f, err := os.Create(".bigpicture.json")
		if err != nil {
			panic(err)
		}
		f.Write([]byte(`{}`))
		return f
	}

	f, err := os.Open(".bigpicture.json")
	if err != nil {
		panic(err)
	}
	return f
}
