package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

const (
	FileName = ".bigpicture.json"
)

type Validator struct {
	Type string         `json:"type"`
	Args map[string]any `json:"args"`
}

type Configuration struct {
	RootDir      string      `json:"root_dir"`
	Lang         string      `json:"lang"`
	IgnoredPaths []string    `json:"ignore"`
	Validators   []Validator `json:"validators"`
	Port         int         `json:"port"`
}

func Init() *Configuration {
	file := checkFileExistAndCreate()
	defer file.Close()

	cfg := &Configuration{
		RootDir: "/",
		Port:    44525,
	}

	err := json.NewDecoder(file).Decode(&cfg)
	if err != nil && !errors.Is(err, io.EOF) {
		panic(err)
	}

	return cfg
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
