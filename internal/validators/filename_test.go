package validators

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFileNameValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})

	args := map[string]any{}
	_, err := NewFileNameValidator(args, nil)
	assert.Equal(t, "module is required and must be string", err.Error())

	args = map[string]any{"module": "a"}
	_, err = NewFileNameValidator(args, nil)
	assert.Equal(t, "max_length is required and must be int", err.Error())

	args = map[string]any{"module": "a", "max_length": "wrong"}
	_, err = NewFileNameValidator(args, nil)
	assert.Equal(t, "max_length is required and must be int", err.Error())

	args = map[string]any{"module": "wrong", "max_length": float64(10)}
	_, err = NewFileNameValidator(args, tree)
	assert.Equal(t, "'wrong' is not a valid module. Path should start with /", err.Error())

	args = map[string]any{"module": "a", "max_length": float64(10)}
	validator, err := NewFileNameValidator(args, tree)
	assert.Nil(t, err)

	args = map[string]any{"module": "a", "max_length": float64(10), "regexp": "[a-z]"}
	validator, err = NewFileNameValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.args.Module)
	assert.Equal(t, 10, validator.args.MaxLength)
	assert.Equal(t, "[a-z]", validator.args.Regexp)
}

func TestFileNameValidator_Validate(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["/srv"] = graph.NewNode("srv", "/srv", graph.Dir, []string{})
	tree.Nodes["/srv/server.go"] = graph.NewNode("srv", "/srv/server.go", graph.Dir, []string{})
	tree.Nodes["/srv/server_1.go"] = graph.NewNode("srv", "/srv/server_1.go", graph.Dir, []string{})

	args := map[string]any{"module": "/srv", "max_length": float64(5)}
	validator, _ := NewFileNameValidator(args, tree)
	err := validator.Validate()
	assert.Equal(t, "File name of '/srv/server.go' is longer than 5", err.Error())

	args = map[string]any{"module": "/srv", "max_length": float64(10)}
	validator, _ = NewFileNameValidator(args, tree)
	err = validator.Validate()
	assert.Nil(t, err)

	args = map[string]any{"module": "/srv", "max_length": float64(10), "regexp": "^[a-z]+$"}
	validator, err = NewFileNameValidator(args, tree)
	err = validator.Validate()
	assert.Equal(t, "File name of '/srv/server_1.go' does not match '^[a-z]+$'", err.Error())

	args = map[string]any{"module": "/srv", "max_length": float64(10), "regexp": "^[a-z]+[0-9_]*$"}
	validator, err = NewFileNameValidator(args, tree)
	err = validator.Validate()
	assert.Nil(t, err)
}
