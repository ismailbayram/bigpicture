package validators

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type SizeValidatorArgs struct {
	Module string   `json:"module" validate:"required=true"`
	Max    float64  `json:"max" validate:"required=true,min=1"`
	Ignore []string `json:"ignore"`
}

type SizeValidator struct {
	args *SizeValidatorArgs
	tree *graph.Tree
}

func NewSizeValidator(args map[string]any, tree *graph.Tree) (*SizeValidator, error) {
	validatorArgs := &SizeValidatorArgs{}
	if err := validateArgs(args, validatorArgs); err != nil {
		return nil, err
	}

	if len(validatorArgs.Module) > 1 && strings.HasSuffix(validatorArgs.Module, "/*") {
		validatorArgs.Module = validatorArgs.Module[:len(validatorArgs.Module)-2]
	}

	module, err := validatePath(validatorArgs.Module, tree)
	if err != nil {
		return nil, err
	}
	validatorArgs.Module = module

	return &SizeValidator{
		args: validatorArgs,
		tree: tree,
	}, nil
}

func (v *SizeValidator) Validate() error {
	totalSize := v.tree.Nodes[v.args.Module].LineCount

	for _, node := range v.tree.Nodes {
		if isIgnored(v.args.Ignore, node.Path) || node.Type == graph.File {
			continue
		}

		nodeSizePercent := float64(node.LineCount) / float64(totalSize) * 100

		if node.Parent == v.args.Module && nodeSizePercent > v.args.Max {
			return fmt.Errorf(
				"Size of module '%s' is %.2f%%, but maximum allowed is %.2f%%",
				node.Path,
				nodeSizePercent,
				v.args.Max,
			)
		}
	}
	return nil
}
