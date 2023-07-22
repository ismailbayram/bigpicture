package validators

import (
	"errors"
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"strings"
)

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
	for _, link := range v.tree.Links {
		if strings.HasPrefix(link.From.Path, v.from) && strings.HasPrefix(link.To.Path, v.to) {
			return errors.New(fmt.Sprintf("'%s' cannot import '%s'", link.From.Path, link.To.Path))
		}
		if v.from == "*" && strings.HasPrefix(link.To.Path, v.to) || v.to == "*" && strings.HasPrefix(link.From.Path, v.from) {
			return errors.New(fmt.Sprintf("'%s' cannot import '%s'", link.From.Path, link.To.Path))
		}
	}

	return nil
}
