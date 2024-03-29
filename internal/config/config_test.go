package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	defer os.Remove(FileName)

	cfg := Init()
	assert.NotNil(t, cfg)
	assert.Equal(t, 44525, cfg.Port)
}

func Test_CheckFileExistAndCreate(t *testing.T) {
	defer os.Remove(FileName)

	f := checkFileExistAndCreate()
	defer f.Close()

	fileExist, _ := os.Stat(FileName)
	assert.Equal(t, FileName, fileExist.Name())
	assert.NotNil(t, f)
}
