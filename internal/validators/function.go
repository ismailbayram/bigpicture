package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type FunctionValidatorArgs struct {
	Module       string   `json:"module" validate:"required=true"`
	MaxLineCount int      `json:"max_line_count" validate:"required=true,gte=1"`
	Ignore       []string `json:"ignore"`
}

type FunctionValidator struct {
	args *FunctionValidatorArgs
	tree *graph.Tree
}

func NewFunctionValidator(args map[string]any, tree *graph.Tree) (*FunctionValidator, error) {
	validatorArgs := &FunctionValidatorArgs{}
	if err := validateArgs(args, validatorArgs); err != nil {
		return nil, err
	}

	if len(validatorArgs.Module) > 1 && strings.HasSuffix(validatorArgs.Module, "/*") {
		validatorArgs.Module = validatorArgs.Module[:len(validatorArgs.Module)-2]
	}

	if err := validatePath(validatorArgs.Module, tree); err != nil {
		return nil, err
	}

	return &FunctionValidator{
		args: validatorArgs,
		tree: tree,
	}, nil
}

func (v *FunctionValidator) Validate() error {
	for _, node := range v.tree.Nodes {
		if isIgnored(v.args.Ignore, node.Path) {
			continue
		}

		if strings.HasPrefix(node.Path, v.args.Module) {
			for _, function := range node.Functions {
				if function.LineCount > v.args.MaxLineCount {
					return errors.New(fmt.Sprintf(
						"Line count of function '%s' in '%s' is %d, but maximum allowed is %d",
						function.Name,
						node.Path,
						function.LineCount,
						v.args.MaxLineCount,
					))
				}
			}
		}
	}

	return nil
}
