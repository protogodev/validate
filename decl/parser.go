package decl

import (
	_ "embed"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//go:embed builtin.go
var BuiltinDecls string

var ErrBadDecl = errors.New("bad declaration of var `_`")

var reVersion = regexp.MustCompile(`(/v[0-9]+)$`)

type Types []string

func (ts Types) Allow(typ types.Type) bool {
	for _, t := range ts {
		switch t {
		case "comparable":
			if IsComparable(typ) {
				return true
			}
		case "ordered":
			if IsOrdered(typ) {
				return true
			}
		case "string":
			if IsString(typ) {
				return true
			}
		case "bytes":
			if IsBytes(typ) {
				return true
			}
		case "slice":
			if IsSlice(typ) {
				return true
			}
		}
	}
	return false
}

type Range struct {
	Min, Max int
}

func (r Range) Contain(n int) bool {
	return n >= r.Min && n <= r.Max
}

type Validator struct {
	Import       string
	Qualifier    string
	Name         string
	IsGeneric    bool
	Alias        string
	AllowedTypes Types
	ArgNum       Range
}

func Parse(decls string) ([]*Validator, error) {
	if decls == "" {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", decls, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	//ast.Print(fset, f)

	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}

		for _, s := range gd.Specs {
			vs, ok := s.(*ast.ValueSpec)
			if !ok || vs.Names[0].Name != "_" {
				continue
			}

			// Here we find the variable declaration of "_".

			p := Parser{
				fset:     fset,
				imports:  parseImports(f.Imports),
				comments: f.Comments,
			}
			return p.Parse(vs.Values)
		}
	}

	return nil, nil
}

func parseImports(imports []*ast.ImportSpec) map[string]string {
	m := make(map[string]string)
	for _, i := range imports {
		path := strings.Trim(i.Path.Value, `"`)
		result := reVersion.FindAllStringSubmatch(path, -1)
		if len(result) > 0 {
			version := result[0][1]
			path = strings.TrimSuffix(path, version)
		}

		name := filepath.Base(path)
		if i.Name != nil {
			name = i.Name.Name
		}

		m[name] = strings.Trim(i.Path.Value, `"`)
	}
	return m
}

type Parser struct {
	fset     *token.FileSet
	imports  map[string]string
	comments []*ast.CommentGroup
}

func (p Parser) Parse(values []ast.Expr) ([]*Validator, error) {
	if len(values) == 0 {
		return nil, nil
	}

	value, ok := values[0].(*ast.CompositeLit)
	if !ok {
		return nil, ErrBadDecl
	}

	var validators []*Validator
	for _, elt := range value.Elts {
		switch e := elt.(type) {
		case *ast.IndexExpr:
			x, ok := e.X.(*ast.SelectorExpr)
			if !ok {
				return nil, ErrBadDecl
			}

			validator, err := p.parseValidator(x)
			if err != nil {
				return nil, err
			}
			validator.IsGeneric = true

			validators = append(validators, validator)

		case *ast.SelectorExpr:
			validator, err := p.parseValidator(e)
			if err != nil {
				return nil, err
			}

			validators = append(validators, validator)

		default:
			return nil, ErrBadDecl
		}
	}

	return validators, nil
}

func (p Parser) parseValidator(e *ast.SelectorExpr) (*Validator, error) {
	qualifier, name, pos, err := p.parseSelector(e)
	if err != nil {
		return nil, err
	}

	comment := p.getComment(pos)
	comment = strings.TrimPrefix(comment, "//")
	if comment == "" {
		return nil, ErrBadDecl
	}

	validator := &Validator{
		Import:    p.imports[qualifier],
		Qualifier: qualifier,
		Name:      name,
		Alias:     strings.ToLower(name), // Defaults to the lowercase version of name.
	}

	fields := strings.Fields(comment)
	for _, f := range fields {
		parts := strings.Split(f, "=")
		k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch k {
		case "name":
			validator.Alias = v
		case "type":
			validator.AllowedTypes = strings.Split(v, "|")
		case "args":
			if strings.HasSuffix(v, "+") {
				v = strings.TrimSuffix(v, "+")
				n := mustAtoi(v)
				validator.ArgNum = Range{Min: n, Max: math.MaxInt}
			} else {
				n := mustAtoi(v)
				validator.ArgNum = Range{Min: n, Max: n}
			}
		}
	}

	return validator, nil
}

func (p Parser) parseSelector(e *ast.SelectorExpr) (qualifier, name string, pos token.Pos, err error) {
	x, ok := e.X.(*ast.Ident)
	if !ok {
		return "", "", 0, ErrBadDecl
	}
	return x.Name, e.Sel.Name, e.Sel.NamePos, nil
}

func (p Parser) getComment(pos token.Pos) string {
	for _, comment := range p.comments {
		c := comment.List[0]
		if p.fset.Position(c.Slash).Line+1 == p.fset.Position(pos).Line {
			return c.Text
		}
	}
	return ""
}

func mustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}
