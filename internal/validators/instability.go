package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type InstabilityValidator struct {
	module string
	max    float64
	tree   *graph.Tree
}

func NewInstabilityValidator(args map[string]any, tree *graph.Tree) (*InstabilityValidator, error) {
	_module, err := validateArg(args, "module", "string")
	if err != nil {
		return nil, err
	}

	_max, err := validateArg(args, "max", "float")
	if err != nil {
		return nil, err
	}

	module := _module.(string)
	max := _max.(float64)
	if max < 0 || max > 1 {
		return nil, errors.New("'max' must be between 0 and 1")
	}

	if len(module) > 1 && strings.HasSuffix(module, "/*") {
		module = module[:len(module)-2]
	}

	if err := validatePath(module, tree); err != nil {
		return nil, err
	}

	return &InstabilityValidator{
		module: module,
		max:    max,
		tree:   tree,
	}, nil
}

func (v *InstabilityValidator) Validate() error {
	node := v.tree.Nodes[v.module]

	if node.Instability != nil && *node.Instability > v.max {
		return errors.New(fmt.Sprintf(
			"instability of %s is %.2f, but should be less than %.2f",
			node.Path,
			*node.Instability,
			v.max,
		))
	}
	return nil
}
