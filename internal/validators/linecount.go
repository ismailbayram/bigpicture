package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type LineCountValidatorArgs struct {
	Module string   `json:"module" validate:"required=true"`
	Max    int      `json:"max" validate:"required=true,min=1"`
	Ignore []string `json:"ignore"`
}

type LineCountValidator struct {
	args *LineCountValidatorArgs
	tree *graph.Tree
}

func NewLineCountValidator(args map[string]any, tree *graph.Tree) (*LineCountValidator, error) {
	validatorArgs := &LineCountValidatorArgs{}
	if err := validateArgs(args, validatorArgs); err != nil {
		return nil, err
	}

	if len(validatorArgs.Module) > 1 && strings.HasSuffix(validatorArgs.Module, "/*") {
		validatorArgs.Module = validatorArgs.Module[:len(validatorArgs.Module)-2]
	}

	if err := validatePath(validatorArgs.Module, tree); err != nil {
		return nil, err
	}

	return &LineCountValidator{
		args: validatorArgs,
		tree: tree,
	}, nil
}

func (v *LineCountValidator) Validate() error {
	for _, node := range v.tree.Nodes {
		if isIgnored(v.args.Ignore, node.Path) {
			continue
		}

		if strings.HasPrefix(node.Path, v.args.Module) && node.LineCount > v.args.Max {
			return errors.New(fmt.Sprintf(
				"Line count of module '%s' is %d, but maximum allowed is %d",
				node.Path,
				node.LineCount,
				v.args.Max,
			))
		}
	}
	return nil
}
