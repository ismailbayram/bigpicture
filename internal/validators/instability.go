package validators

import "github.com/ismailbayram/bigpicture/internal/graph"

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
