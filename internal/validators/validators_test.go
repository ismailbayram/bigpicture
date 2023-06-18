package validators

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})

	validator, err := NewValidator("no_import", map[string]any{"from": "a", "to": "b"}, tree)
	assert.Nil(t, err)
	assert.IsType(t, &NoImportValidator{}, validator)

	validator, err = NewValidator("instability", map[string]any{"module": "a", "max": 0.5}, tree)
	assert.Nil(t, err)
	assert.IsType(t, &InstabilityValidator{}, validator)

	validator, err = NewValidator("unknown", map[string]any{}, tree)
	assert.NotNil(t, err)
	assert.Nil(t, validator)
	assert.Equal(t, "unknown validator type: unknown", err.Error())
}

func TestNewNoImportValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})

	args := map[string]any{}
	_, err := NewNoImportValidator(args, nil)
	assert.Equal(t, "'from' is required", err.Error())

	args = map[string]any{"from": "a"}
	_, err = NewNoImportValidator(args, nil)
	assert.Equal(t, "'to' is required", err.Error())

	args = map[string]any{"from": ""}
	_, err = NewNoImportValidator(args, nil)
	assert.Equal(t, "'from' cannot be empty", err.Error())

	args = map[string]any{"from": "a", "to": ""}
	_, err = NewNoImportValidator(args, nil)
	assert.Equal(t, "'to' cannot be empty", err.Error())

	args = map[string]any{"from": "wrong", "to": "b/*"}
	validator, err := NewNoImportValidator(args, tree)
	assert.Equal(t, "'wrong' is not a valid module. Path should start with /", err.Error())

	args = map[string]any{"from": "a", "to": "*"}
	validator, err = NewNoImportValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.from)
	assert.Equal(t, "*", validator.to)
	assert.NotNil(t, validator.tree)

	args = map[string]any{"from": "a/*", "to": "b/*"}
	validator, err = NewNoImportValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.from)
	assert.Equal(t, "b", validator.to)

	args = map[string]any{"from": "a/*", "to": "b/*"}
	validator, err = NewNoImportValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.from)
	assert.Equal(t, "b", validator.to)

}

func TestNoImportValidator_Validate(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["server"] = graph.NewNode("server", "server", graph.Dir, []string{
		"browser/go",
	})
	tree.Nodes["config"] = graph.NewNode("config", "config", graph.Dir, []string{})
	tree.Nodes["config/subconfig"] = graph.NewNode("subconfig", "config/subconfig", graph.Dir, []string{})
	tree.Nodes["browser"] = graph.NewNode("browser", "browser", graph.Dir, []string{
		"config/subconfig",
	})
	tree.Nodes["browser/go"] = graph.NewNode("go", "browser/go", graph.Dir, []string{
		"config/subconfig",
	})
	tree.GenerateLinks()

	validator, err := NewValidator("no_import", map[string]any{"from": "server", "to": "config/subconfig"}, tree)
	assert.Nil(t, err)
	assert.Nil(t, validator.Validate())

	validator, err = NewValidator("no_import", map[string]any{"from": "browser", "to": "config"}, tree)
	assert.Nil(t, err)
	assert.NotNil(t, validator.Validate())
	assert.Equal(t, "'browser' cannot import 'config/subconfig'", validator.Validate().Error())

	validator, err = NewValidator("no_import", map[string]any{"from": "*", "to": "browser"}, tree)
	assert.Nil(t, err)
	assert.NotNil(t, validator.Validate())
	assert.Equal(t, "'server' cannot import 'browser/go'", validator.Validate().Error())
}
