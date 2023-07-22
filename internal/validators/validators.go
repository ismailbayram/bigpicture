package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
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
	case "line_count":
		return NewLineCountValidator(args, tree)
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
	case "int":
		_, ok := val.(float64)
		if !ok {
			return nil, errors.New(fmt.Sprintf("'%s' must be an integer", arg))
		}
		return int(val.(float64)), nil
	}

	return val, nil
}

func validatePath(path string, tree *graph.Tree) error {
	if _, ok := tree.Nodes[path]; !ok && path != "*" {
		return errors.New(fmt.Sprintf("'%s' is not a valid module. Path should start with /", path))
	}
	return nil
}
