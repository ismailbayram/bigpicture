package validators

import (
	"github.com/ismailbayram/bigpicture/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFunctionValidator(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["a"] = graph.NewNode("a", "a", graph.Dir, []string{})
	tree.Nodes["b"] = graph.NewNode("b", "b", graph.Dir, []string{})

	args := map[string]any{}
	_, err := NewFunctionValidator(args, nil)
	assert.Equal(t, "module is required and must be string", err.Error())

	args = map[string]any{"module": "a"}
	_, err = NewFunctionValidator(args, nil)
	assert.Equal(t, "max_line_count is required and must be int", err.Error())

	args = map[string]any{"module": "a", "max_line_count": "wrong"}
	_, err = NewFunctionValidator(args, nil)
	assert.Equal(t, "max_line_count is required and must be int", err.Error())

	args = map[string]any{"module": "wrong", "max_line_count": float64(100)}
	_, err = NewFunctionValidator(args, tree)
	assert.Equal(t, "'wrong' is not a valid module. Path should start with /", err.Error())

	args = map[string]any{"module": "a", "max_line_count": float64(100)}
	validator, err := NewFunctionValidator(args, tree)
	assert.Nil(t, err)
	assert.Equal(t, "a", validator.args.Module)
	assert.Equal(t, 100, validator.args.MaxLineCount)
}

func TestFunctionValidator_Validate(t *testing.T) {
	tree := graph.NewTree("root")
	tree.Nodes["/server"] = graph.NewNode("server", "/server", graph.Dir, []string{})
	tree.Nodes["/server/server.go"] = graph.NewNode("server", "/server/server.go", graph.Dir, []string{})
	tree.Nodes["/server/server.go"].Functions = []graph.Function{
		{Name: "func1", LineCount: 150},
	}

	args := map[string]any{"module": "/server", "max_line_count": float64(100)}
	validator, _ := NewFunctionValidator(args, tree)
	err := validator.Validate()
	assert.Equal(t, "Line count of function 'func1' in '/server/server.go' is 150, but maximum allowed is 100", err.Error())

	args = map[string]any{"module": "/server", "max_line_count": float64(150)}
	validator, _ = NewFunctionValidator(args, tree)
	err = validator.Validate()
	assert.Nil(t, err)
}
