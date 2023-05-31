package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

const (
	FileName = ".bigpicture.json"
)

type Configuration struct {
	IgnoredPaths []string `json:"ignore"`
	Validators   []string `json:"validators"`
	Port         int      `json:"port"`
}

func Init() *Configuration {
	file := checkFileExistAndCreate()
	defer file.Close()

	cfg := &Configuration{
		Port: 44525,
	}

	err := json.NewDecoder(file).Decode(&cfg)
	if err != nil && !errors.Is(err, io.EOF) {
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
	_, err := os.Stat(FileName)
	if err != nil {
		return createFile(err)
	}

	f, err := os.Open(FileName)
	if err != nil {
		panic(err)
	}
	return f
}

func createFile(err error) *os.File {
	if !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	f, err := os.Create(FileName)
	f.Write([]byte(`{}`))
	if err != nil {
		panic(err)
	}
	return f
}
