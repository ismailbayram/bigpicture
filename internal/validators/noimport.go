package validators

import (
	"errors"
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

	if len(validatorArgs.From) > 1 && strings.HasSuffix(validatorArgs.From, "/*") {
		validatorArgs.From = validatorArgs.From[:len(validatorArgs.From)-2]
	}

	if len(validatorArgs.To) > 1 && strings.HasSuffix(validatorArgs.To, "/*") {
		validatorArgs.To = validatorArgs.To[:len(validatorArgs.To)-2]
	}
	if err := validatePath(validatorArgs.From, tree); err != nil {
		return nil, err
	}

	if err := validatePath(validatorArgs.To, tree); err != nil {
		return nil, err
	}

	return &NoImportValidator{
		args: validatorArgs,
		tree: tree,
	}, nil
}

func (v *NoImportValidator) Validate() error {
	for _, link := range v.tree.Links {
		if strings.HasPrefix(link.From.Path, v.args.From) && strings.HasPrefix(link.To.Path, v.args.To) {
			return errors.New(fmt.Sprintf("'%s' cannot import '%s'", link.From.Path, link.To.Path))
		}
		if v.args.From == "*" && strings.HasPrefix(link.To.Path, v.args.To) || v.args.To == "*" && strings.HasPrefix(link.From.Path, v.args.From) {
			return errors.New(fmt.Sprintf("'%s' cannot import '%s'", link.From.Path, link.To.Path))
		}
	}

	return nil
}
