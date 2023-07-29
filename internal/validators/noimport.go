package validators

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type NoImportValidatorArgs struct {
	From string `json:"from" validate:"required=true"`
	To   string `json:"to" validate:"required=true"`
}

type NoImportValidator struct {
	args *NoImportValidatorArgs
	tree *graph.Tree
}

func NewNoImportValidator(args map[string]any, tree *graph.Tree) (*NoImportValidator, error) {
	validatorArgs := &NoImportValidatorArgs{}
	if err := validateArgs(args, validatorArgs); err != nil {
		return nil, err
	}

	from, err := validatePath(validatorArgs.From, tree)
	if err != nil {
		return nil, err
	}
	validatorArgs.From = from

	to, err := validatePath(validatorArgs.To, tree)
	if err != nil {
		return nil, err
	}
	validatorArgs.To = to

	return &NoImportValidator{
		args: validatorArgs,
		tree: tree,
	}, nil
}

func (v *NoImportValidator) Validate() error {
	for _, link := range v.tree.Links {
		if strings.HasPrefix(link.From.Path, v.args.From) && strings.HasPrefix(link.To.Path, v.args.To) {
			return fmt.Errorf("'%s' cannot import '%s'", link.From.Path, link.To.Path)
		}
		if v.args.From == "*" && strings.HasPrefix(link.To.Path, v.args.To) || v.args.To == "*" && strings.HasPrefix(link.From.Path, v.args.From) {
			return fmt.Errorf("'%s' cannot import '%s'", link.From.Path, link.To.Path)
		}
	}

	return nil
}
