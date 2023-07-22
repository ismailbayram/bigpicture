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

	validator, err = NewValidator("line_count", map[string]any{"module": "a", "max": float64(100)}, tree)
	assert.Nil(t, err)
	assert.IsType(t, &LineCountValidator{}, validator)

	validator, err = NewValidator("unknown", map[string]any{}, tree)
	assert.NotNil(t, err)
	assert.Nil(t, validator)
	assert.Equal(t, "unknown validator type: unknown", err.Error())
}
