package validators

import (
	"encoding/json"
	"errors"
	"fmt"
	v10Validator "github.com/go-playground/validator/v10"
	"github.com/ismailbayram/bigpicture/internal/graph"
	"regexp"
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
	case "line_count":
		return NewLineCountValidator(args, tree)
	case "function":
		return NewFunctionValidator(args, tree)
	default:
		return nil, errors.New(fmt.Sprintf("unknown validator type: %s", t))
	}
}

func validateArgs(args map[string]any, validatorArgStruct any) error {
	jsonData, _ := json.Marshal(args)
	_ = json.Unmarshal(jsonData, validatorArgStruct)
	validate := v10Validator.New()
	err := validate.Struct(validatorArgStruct)
	if err != nil {
		return validationErrorToText(err.(v10Validator.ValidationErrors)[0])
	}
	return nil
}

func validationErrorToText(e v10Validator.FieldError) error {
	word := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Errorf("%s is required and must be %s", word, e.Type().String())
	case "max":
		return fmt.Errorf("%s cannot be longer than %s", word, e.Param())
	case "min":
		return fmt.Errorf("%s must be longer than %s", word, e.Param())
	case "gte":
		return fmt.Errorf("%s must be greater than or equal to %s", word, e.Param())
	case "lte":
		return fmt.Errorf("%s must be less than or equal to %s", word, e.Param())
	case "email":
		return fmt.Errorf("%s is not a valid email address", word)
	case "len":
		return fmt.Errorf("%s must be %s characters long", word, e.Param())
	}
	return fmt.Errorf("%s is not %s", word, e.Type().String())
}

func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func validateArg(args map[string]any, arg string, argType string, required bool) (any, error) {
	val, ok := args[arg]
	if !ok && required {
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
	case "float":
		_, ok := val.(float64)
		if !ok {
			return nil, errors.New(fmt.Sprintf("'%s' must be a float", arg))
		}
		return val.(float64), nil
	}

	return val, nil
}

func validatePath(path string, tree *graph.Tree) error {
	if _, ok := tree.Nodes[path]; !ok && path != "*" {
		return errors.New(fmt.Sprintf("'%s' is not a valid module. Path should start with /", path))
	}
	return nil
}

func isIgnored(ignoreList []string, path string) bool {
	isIgnored := false
	for _, ignore := range ignoreList {
		regxp := ignore
		if strings.HasPrefix(ignore, "*") {
			regxp = fmt.Sprintf("^%s$", ignore)
		}
		re := regexp.MustCompile(regxp)
		if re.MatchString(path) {
			isIgnored = true
			break
		}
	}
	return isIgnored
}
