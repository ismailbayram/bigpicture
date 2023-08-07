package validators

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLineCountValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})
	tree.Nodes["a"].LineCount = 100
	tree.Nodes["b"].LineCount = 200

	args := map[string]any{}
	_, err := NewLineCountValidator(args, nil)
	assert.Equal(t, "module is required and must be string", err.Error())

	args = map[string]any{"module": "a"}
	_, err = NewLineCountValidator(args, nil)
	assert.Equal(t, "max is required and must be int", err.Error())

	args = map[string]any{"module": "a", "max": "wrong"}
	_, err = NewLineCountValidator(args, nil)
	assert.Equal(t, "max is required and must be int", err.Error())

	args = map[string]any{"module": "wrong", "max": float64(100)}
	_, err = NewLineCountValidator(args, tree)
	assert.Equal(t, "'wrong' is not a valid module. Path should start with /", err.Error())

	args = map[string]any{"module": "a", "max": float64(100)}
	validator, err := NewLineCountValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.args.Module)
	assert.Equal(t, 100, validator.args.Max)
}

func TestLineCountValidator_Validate(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})
	tree.Nodes["a"].LineCount = 100
	tree.Nodes["b"].LineCount = 200

	args := map[string]any{"module": "a", "max": float64(100)}
	validator, _ := NewLineCountValidator(args, tree)
	err := validator.Validate()
	assert.Nil(t, err)

	args = map[string]any{"module": "a", "max": float64(50)}
	validator, _ = NewLineCountValidator(args, tree)
	err = validator.Validate()
	assert.Equal(t, "Line count of module 'a' is 100, but maximum allowed is 50", err.Error())

	tree = graph.NewTree("root")
	tree.Nodes["server"] = graph.NewNode("server", "server", graph.Dir, []string{
		"browser/go",
	})
	tree.Nodes["config"] = graph.NewNode("config", "config", graph.Dir, []string{})
	tree.Nodes["config/subconfig"] = graph.NewNode("subconfig", "config/subconfig", graph.Dir, []string{})
	tree.Nodes["config/subconfig"].LineCount = 200
	tree.Nodes["browser"] = graph.NewNode("browser", "browser", graph.Dir, []string{
		"config/subconfig",
	})
	tree.Nodes["browser/go"] = graph.NewNode("go", "browser/go", graph.Dir, []string{
		"config/subconfig",
	})

	args = map[string]any{"module": "config", "max": float64(100)}
	validator, _ = NewLineCountValidator(args, tree)
	err = validator.Validate()
	assert.Equal(t, "Line count of module 'config/subconfig' is 200, but maximum allowed is 100", err.Error())
}
