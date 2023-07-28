package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
)

type InstabilityValidatorArgs struct {
	Module string  `json:"module" validate:"required=true"`
	Max    float64 `json:"max" validate:"required=true,min=0,max=1"`
}

type InstabilityValidator struct {
	args *InstabilityValidatorArgs
	tree *graph.Tree
}

func NewInstabilityValidator(args map[string]any, tree *graph.Tree) (*InstabilityValidator, error) {
	validatorArgs := &InstabilityValidatorArgs{}
	if err := validateArgs(args, validatorArgs); err != nil {
		return nil, err
	}

	module, err := validatePath(validatorArgs.Module, tree)
	if err != nil {
		return nil, err
	}
	validatorArgs.Module = module

	return &InstabilityValidator{
		args: validatorArgs,
		tree: tree,
	}, nil
}

func (v *InstabilityValidator) Validate() error {
	node := v.tree.Nodes[v.args.Module]

	if node.Instability != nil && *node.Instability > v.args.Max {
		return errors.New(fmt.Sprintf(
			"instability of %s is %.2f, but should be less than %.2f",
			node.Path,
			*node.Instability,
			v.args.Max,
		))
	}
	return nil
}
