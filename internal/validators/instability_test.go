package validators

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInstabilityValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})

	args := map[string]any{}
	_, err := NewInstabilityValidator(args, nil)
	assert.Equal(t, "'module' is required", err.Error())

	args = map[string]any{"module": "a"}
	_, err = NewInstabilityValidator(args, nil)
	assert.Equal(t, "'max' is required", err.Error())

	args = map[string]any{"module": "a", "max": "wrong"}
	_, err = NewInstabilityValidator(args, nil)
	assert.Equal(t, "'max' must be a float", err.Error())

	args = map[string]any{"module": "wrong", "max": float64(100)}
	_, err = NewInstabilityValidator(args, tree)
	assert.Equal(t, "'max' must be between 0 and 1", err.Error())

	args = map[string]any{"module": "wrong", "max": 0.2}
	_, err = NewInstabilityValidator(args, tree)
	assert.Equal(t, "'wrong' is not a valid module. Path should start with /", err.Error())

	args = map[string]any{"module": "a", "max": 0.2}
	validator, err := NewInstabilityValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.module)
	assert.Equal(t, 0.2, validator.max)
}

func TestInstabilityValidator_Validate(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["/server"] = graph.NewNode("server", "/server", graph.Dir, []string{})
	tree.Nodes["/server/server.go"] = graph.NewNode("server", "/server/server.go", graph.Dir, []string{
		"/browser/go.go",
	})
	tree.Nodes["/config"] = graph.NewNode("config", "/config", graph.Dir, []string{})
	tree.Nodes["/browser"] = graph.NewNode("browser", "/browser", graph.Dir, []string{})
	tree.Nodes["/browser/go.go"] = graph.NewNode("browser", "/browser/go.go", graph.Dir, []string{
		"/config",
		"/graph/graph.go",
	})
	tree.Nodes["/graph"] = graph.NewNode("graph", "/graph", graph.Dir, []string{})
	tree.Nodes["/graph/graph.go"] = graph.NewNode("graph", "/graph/graph.go", graph.Dir, []string{})
	tree.GenerateLinks()
	tree.CalculateInstability()

	args := map[string]any{"module": "/browser", "max": 0.25}
	validator, _ := NewInstabilityValidator(args, tree)
	err := validator.Validate()
	assert.Equal(t, "instability of /browser is 0.50, but should be less than 0.25", err.Error())

	args = map[string]any{"module": "/browser", "max": 0.50}
	validator, _ = NewInstabilityValidator(args, tree)
	err = validator.Validate()
	assert.Nil(t, err)
}
