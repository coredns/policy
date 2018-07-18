package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type errors struct {
	src     string
	Pkg     string        `yaml:"package"`
	Imports []string      `yaml:"import"`
	Errors  []*definition `yaml:"errors"`
}

type definition struct {
	ID     string      `yaml:"id"`
	Fields []*field    `yaml:"fields"`
	Msg    string      `yaml:"msg"`
	Args   []*argument `yaml:"args"`
	Desc   string      `yaml:"desc"`
}

type field struct {
	ID   string `yaml:"id"`
	Type string `yaml:"type"`
}

type argument struct {
	Field   string   `yaml:"field"`
	Expr    string   `yaml:"expr"`
	Snippet *snippet `yaml:"snippet"`
}

type snippet struct {
	Result string `yaml:"result"`
	Code   string `yaml:"code"`
}

func unmarshal(path string) (*errors, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	e := &errors{
		src: filepath.Base(path)}
	err = yaml.Unmarshal(b, &e)
	if err != nil {
		return nil, err
	}

	if len(e.Errors) <= 0 {
		return nil, fmt.Errorf("No errors defined")
	}

	e.Pkg = strings.TrimSpace(e.Pkg)
	if len(e.Pkg) <= 0 {
		return nil, fmt.Errorf("Missing package name")
	}

	for i, item := range e.Errors {
		ID := item.ID

		item.ID = strings.TrimSpace(ID)
		if len(item.ID) <= 0 {
			return nil, fmt.Errorf("Empty error id for error at %d", i+1)
		}

		if strings.ToUpper(item.ID[:1]) == item.ID[:1] && len(item.Desc) <= 0 {
			return nil, fmt.Errorf("Exported error %q has no description", item.ID)
		}

		for j, field := range item.Fields {
			fID := field.ID

			field.ID = strings.TrimSpace(fID)
			if len(field.ID) <= 0 {
				return nil, fmt.Errorf("Empty field id at %d for error %q", j+1, ID)
			}

			field.Type = strings.TrimSpace(field.Type)
			if len(field.Type) <= 0 {
				return nil, fmt.Errorf("Empty %q field type for error %q", fID, ID)
			}
		}

		if len(item.Msg) <= 0 {
			return nil, fmt.Errorf("Empty error message for error %q", ID)
		}

		for j, arg := range item.Args {
			arg.Field = strings.TrimSpace(arg.Field)
			field := len(arg.Field) > 0

			arg.Expr = strings.TrimSpace(arg.Expr)
			expr := len(arg.Expr) > 0

			snippet := arg.Snippet != nil

			if !field && !expr && !snippet {
				return nil, fmt.Errorf("Empty argument %d definition for error %q", j+1, ID)
			}

			if field && expr && snippet {
				return nil, fmt.Errorf("Ambiguous argument %d definition for error %q: "+
					"all field, expression and snippet are defined", j+1, ID)
			}

			if field && expr {
				return nil, fmt.Errorf("Ambiguous argument %d definition for error %q: "+
					"both field and expression are defined", j+1, ID)
			}

			if field && snippet {
				return nil, fmt.Errorf("Ambiguous argument %d definition for error %q: "+
					"both field and snippet are defined", j+1, ID)
			}

			if expr && snippet {
				return nil, fmt.Errorf("Ambiguous argument %d definition for error %q: "+
					"both expression and snippet are defined", j+1, ID)
			}

			if snippet {
				arg.Snippet.Result = strings.TrimSpace(arg.Snippet.Result)
				if len(arg.Snippet.Result) <= 0 {
					return nil, fmt.Errorf("Empty result name of snippet argument %d for error %q", j+1, ID)
				}

				code := strings.TrimSpace(arg.Snippet.Code)
				if len(code) <= 0 {
					return nil, fmt.Errorf("Empty snippet argument %d for error %q", j+1, ID)
				}
			}
		}
	}

	imps := make([]string, len(e.Imports))
	for i, imp := range e.Imports {
		imp = strings.TrimSpace(imp)
		if len(imp) <= 0 {
			return nil, fmt.Errorf("Empty import %d", i+1)
		}

		imps[i] = fmt.Sprintf("%q", imp)
	}

	e.Imports = imps
	return e, nil
}
