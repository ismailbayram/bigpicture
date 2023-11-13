package validators

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSizeValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})
	tree.Nodes["a"].LineCount = 100
	tree.Nodes["b"].LineCount = 200

	args := map[string]any{}
	_, err := NewSizeValidator(args, nil)
	assert.Equal(t, "module is required and must be string", err.Error())

	args = map[string]any{"module": "a"}
	_, err = NewSizeValidator(args, nil)
	assert.Equal(t, "max is required and must be float64", err.Error())

	args = map[string]any{"module": "a", "max": "wrong"}
	_, err = NewSizeValidator(args, nil)
	assert.Equal(t, "max is required and must be float64", err.Error())

	args = map[string]any{"module": "wrong", "max": float64(100)}
	_, err = NewSizeValidator(args, tree)
	assert.Equal(t, "'wrong' is not a valid module. Path should start with /", err.Error())

	args = map[string]any{"module": "a", "max": float64(100)}
	validator, err := NewSizeValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.args.Module)
	assert.Equal(t, float64(100), validator.args.Max)
}

func TestSizeValidator_Validate(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["."].LineCount = 600
	tree.Nodes["/dir1"] = graph.NewNode("dir1", "/dir1", graph.Dir, []string{})
	tree.Nodes["/dir1"].LineCount = 100
	tree.Nodes["/dir2"] = graph.NewNode("dir2", "/dir2", graph.Dir, []string{})
	tree.Nodes["/dir2"].LineCount = 200
	tree.Nodes["/dir3"] = graph.NewNode("dir3", "/dir3", graph.Dir, []string{})
	tree.Nodes["/dir3"].LineCount = 300

	args := map[string]any{"module": ".", "max": float64(40)}
	validator, _ := NewSizeValidator(args, tree)
	err := validator.Validate()
	assert.Equal(t, "Size of module '/dir3' is 50.00%, but maximum allowed is 40.00%", err.Error())

	args = map[string]any{"module": ".", "max": float64(50)}
	validator, _ = NewSizeValidator(args, tree)
	err = validator.Validate()
	assert.Nil(t, err)
}
