package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type FunctionValidator struct {
	module       string
	maxLineCount int
	tree         *graph.Tree
}

func NewFunctionValidator(args map[string]any, tree *graph.Tree) (*FunctionValidator, error) {
	_module, err := validateArg(args, "module", "string")
	if err != nil {
		return nil, err
	}

	_max, err := validateArg(args, "max_line_count", "int")
	if err != nil {
		return nil, err
	}

	module := _module.(string)
	max := _max.(int)

	if len(module) > 1 && strings.HasSuffix(module, "/*") {
		module = module[:len(module)-2]
	}

	if err := validatePath(module, tree); err != nil {
		return nil, err
	}

	return &FunctionValidator{
		module:       module,
		maxLineCount: max,
		tree:         tree,
	}, nil
}

func (v *FunctionValidator) Validate() error {
	for _, node := range v.tree.Nodes {
		if strings.HasPrefix(node.Path, v.module) {
			for _, function := range node.Functions {
				if function.LineCount > v.maxLineCount {
					return errors.New(fmt.Sprintf(
						"Line count of function '%s' in '%s' is %d, but maximum allowed is %d",
						function.Name,
						node.Path,
						function.LineCount,
						v.maxLineCount,
					))
				}
			}
		}
	}

	return nil
}
