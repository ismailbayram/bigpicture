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
	for _, node := range v.tree.Nodes {
		if node.Parent != v.module {
			continue
		}

		importCount := 0   // from node to other modules
		importedCount := 0 // from other modules to node
		for _, link := range v.tree.Links {
			if link.From.Path == node.Path && link.To.Parent == node.Parent {
				importCount += 1
			}
			if strings.HasPrefix(link.To.Path, node.Path) && link.From.Parent == node.Parent {
				importedCount += 1
			}
		}

		instability := float64(importCount) / float64(importedCount+importCount)
		if instability > v.max {
			return errors.New(fmt.Sprintf(
				"instability of %s is %.2f, but should be less than %.2f",
				node.Path,
				instability,
				v.max,
			))
		}
	}
	return nil
}
