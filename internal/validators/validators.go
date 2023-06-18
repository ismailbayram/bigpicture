package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

type Validator interface {
	Validate() error
}

func NewValidator(t string, args map[string]any, tree *graph.Tree) (Validator, error) {
	switch t {
	case "no_import":
		return NewNoImportValidator(args, tree)
	case "instability":
		return NewInstabilityValidator(args, tree)
	default:
		return nil, errors.New(fmt.Sprintf("unknown validator type: %s", t))
	}
}

func validateArg(args map[string]any, arg string, argType string) (any, error) {
	val, ok := args[arg]
	if !ok {
		return nil, errors.New(fmt.Sprintf("'%s' is required", arg))
	}

	switch argType {
	case "string":
		_, ok := val.(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("'%s' must be a string", arg))
		}
		if val == "" {
			return nil, errors.New(fmt.Sprintf("'%s' cannot be empty", arg))
		}
		return val, nil
	}
	return val, nil
}

func validatePath(path string, tree *graph.Tree) error {
	if _, ok := tree.Nodes[path]; !ok && path != "*" {
		return errors.New(fmt.Sprintf("'%s' is not a valid module. Path should start with /", path))
	}
	return nil
}

type NoImportValidator struct {
	from string
	to   string
	tree *graph.Tree
}

func NewNoImportValidator(args map[string]any, tree *graph.Tree) (*NoImportValidator, error) {
	_from, err := validateArg(args, "from", "string")
	if err != nil {
		return nil, err
	}
	_to, err := validateArg(args, "to", "string")
	if err != nil {
		return nil, err
	}

	from := _from.(string)
	to := _to.(string)

	if len(from) > 1 && strings.HasSuffix(from, "/*") {
		from = from[:len(from)-2]
	}

	if len(to) > 1 && strings.HasSuffix(to, "/*") {
		to = to[:len(to)-2]
	}
	if err := validatePath(from, tree); err != nil {
		return nil, err
	}

	if err := validatePath(to, tree); err != nil {
		return nil, err
	}

	return &NoImportValidator{
		from: from,
		to:   to,
		tree: tree,
	}, nil
}

func (v *NoImportValidator) Validate() error {
	return nil
}

type InstabilityValidator struct {
	module string
	max    float32
	tree   *graph.Tree
}

func NewInstabilityValidator(args map[string]any, tree *graph.Tree) (*InstabilityValidator, error) {
	return nil, nil
}

func (v *InstabilityValidator) Validate() error {
	return nil
}
