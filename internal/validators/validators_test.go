package validators

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = &graph.Node{PackageName: "a"}
	tree.Nodes["b"] = &graph.Node{PackageName: "b"}

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
	tree.Nodes["a"] = &graph.Node{PackageName: "a"}
	tree.Nodes["b"] = &graph.Node{PackageName: "b"}

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

}
