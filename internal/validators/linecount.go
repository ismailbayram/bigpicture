package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type LineCountValidator struct {
	module string
	max    int
	tree   *graph.Tree
}

func NewLineCountValidator(args map[string]any, tree *graph.Tree) (*LineCountValidator, error) {
	_module, err := validateArg(args, "module", "string")
	if err != nil {
		return nil, err
	}

	_max, err := validateArg(args, "max", "int")
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

	return &LineCountValidator{
		module: module,
		max:    max,
		tree:   tree,
	}, nil
}

func (v *LineCountValidator) Validate() error {
	for _, node := range v.tree.Nodes {
		if strings.HasPrefix(node.Path, v.module) && node.LineCount > v.max {
			return errors.New(fmt.Sprintf(
				"Line count of module '%s' is %d, but maximum allowed is %d",
				node.Path,
				node.LineCount,
				v.max,
			))
		}
	}
	return nil
}
