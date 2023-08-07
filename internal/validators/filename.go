package validators

import (
	"fmt"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"regexp"
	"strings"
)

type FileNameValidatorArgs struct {
	Module    string   `json:"module" validate:"required=true"`
	MaxLength int      `json:"max_length" validate:"required=true,gte=1"`
	Regexp    string   `json:"regexp"`
	Ignore    []string `json:"ignore"`
}

type FileNameValidator struct {
	args *FileNameValidatorArgs
	tree *graph.Tree
}

func NewFileNameValidator(args map[string]any, tree *graph.Tree) (*FileNameValidator, error) {
	validatorArgs := &FileNameValidatorArgs{}
	if err := validateArgs(args, validatorArgs); err != nil {
		return nil, err
	}

	module, err := validatePath(validatorArgs.Module, tree)
	if err != nil {
		return nil, err
	}
	validatorArgs.Module = module

	return &FileNameValidator{
		args: validatorArgs,
		tree: tree,
	}, nil
}

func (v *FileNameValidator) Validate() error {
	for _, node := range v.tree.Nodes {
		if isIgnored(v.args.Ignore, node.Path) {
			continue
		}
		if !strings.HasPrefix(node.Path, v.args.Module) {
			continue
		}

		if len(node.FileName()) > v.args.MaxLength {
			return fmt.Errorf("File name of '%s' is longer than %d", node.Path, v.args.MaxLength)
		}

		if v.args.Regexp != "" {
			re := regexp.MustCompile(v.args.Regexp)
			if !re.MatchString(node.FileName()) {
				return fmt.Errorf("File name of '%s' does not match '%s'", node.Path, v.args.Regexp)
			}
		}
	}

	return nil
}
