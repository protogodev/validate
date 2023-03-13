package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"github.com/protogodev/validate/decl"
)

const (
	DefaultQualifier = "v"
)

type Param struct {
	Name string
	Type types.Type
}

type Validator interface {
	Bind(param Param, decls map[string][]*decl.Validator) error

	// ExprString returns the validating-style expression string.
	ExprString() string
}

// LeafValidator is an expression that represents a leaf validator.
type LeafValidator struct {
	Name string
	Args []string
	Msg  string

	Param Param
	Decls []*decl.Validator
}

func (v *LeafValidator) Bind(param Param, decls map[string][]*decl.Validator) error {
	v.Param = param
	v.Decls = decls[v.Name]

	return v.validate()
}

func (v *LeafValidator) ExprString() string {
	qualifiedName := v.buildQualifiedName()

	args := strings.Join(v.Args, ", ")
	if v.Name == "match" {
		args = fmt.Sprintf("regexp.MustCompile(%s)", args)
	}

	if v.Msg == "" {
		return fmt.Sprintf("%s(%s)", qualifiedName, args)
	}
	return fmt.Sprintf("%s(%s).Msg(%s)", qualifiedName, args, v.Msg)
}

func (v *LeafValidator) validate() error {
	// Special case for validator `_`.
	if v.Name == "_" {
		if !decl.IsStruct(v.Param.Type) {
			return fmt.Errorf("cannot use validator `%s` on type %T", v.Name, v.Param.Type.Underlying())
		}
		return nil
	}

	if len(v.Decls) == 0 {
		return fmt.Errorf("unrecognized validator %q", v.Name)
	}

	// Try to find the first matched declaration.
	idx := -1
	for i, d := range v.Decls {
		if d.AllowedTypes.Allow(v.Param.Type) {
			idx = i
			break
		}
	}
	if idx == -1 {
		// Found no match, return an error.
		return fmt.Errorf("cannot use validator `%s` on type %T", v.Name, v.Param.Type.Underlying())
	}

	// Apply the argument number constraint from the above matched declaration.
	d := v.Decls[idx]
	if !d.ArgNum.Contain(len(v.Args)) {
		return fmt.Errorf("wrong number of arguments for validator %q", v.Name)
	}

	return nil
}

func (v *LeafValidator) buildQualifiedName() string {
	// Special case for validator `_`.
	if v.Name == "_" {
		return v.Param.Name + ".Schema"
	}

	for _, d := range v.Decls {
		if d.AllowedTypes.Allow(v.Param.Type) {
			// Return the qualified name of the first matched declaration.

			name := d.Qualifier + "." + d.Name
			if d.IsGeneric {
				name += "[" + v.Param.Type.String() + "]"
			}
			return name
		}
	}

	return ""
}

// LogicValidator is an expression that represents a logic validator (i.e. `Not`, `And/All` or `Or/Any`).
type LogicValidator struct {
	Qualifier string
	Name      string
	Left      Validator
	Right     Validator // `Not` has no Right validator.
}

func (v *LogicValidator) Bind(param Param, decls map[string][]*decl.Validator) error {
	if err := v.Left.Bind(param, decls); err != nil {
		return err
	}
	if v.Right != nil {
		return v.Right.Bind(param, decls)
	}
	return nil
}

func (v *LogicValidator) ExprString() string {
	qualifiedName := v.buildQualifiedName()

	if v.Right != nil {
		return fmt.Sprintf("%s(%s, %s)", qualifiedName, v.Left.ExprString(), v.Right.ExprString())
	}
	return fmt.Sprintf("%s(%s)", qualifiedName, v.Left.ExprString())
}

func (v *LogicValidator) buildQualifiedName() string {
	switch v.Name {
	case "!":
		return v.Qualifier + ".Not"
	case "&&":
		return v.Qualifier + ".All"
	case "||":
		return v.Qualifier + ".Any"
	}

	return ""
}

func Parse(s string) (Validator, error) {
	expr, err := parser.ParseExpr(s)
	if err != nil {
		return nil, err
	}
	//ast.Print(token.NewFileSet(), expr)

	v, err := Parser{S: s}.Parse(expr)
	if err != nil {
		return nil, err
	}

	return v, nil
}

type Parser struct {
	S string
}

func (p Parser) Parse(e ast.Expr) (Validator, error) {
	switch expr := e.(type) {
	case *ast.UnaryExpr:
		switch expr.Op {
		case token.NOT:
			x, err := p.Parse(expr.X)
			if err != nil {
				return nil, err
			}
			return &LogicValidator{
				Qualifier: DefaultQualifier,
				Name:      "!",
				Left:      x,
			}, nil
		}

	case *ast.BinaryExpr:
		switch expr.Op {
		case token.LAND:
			// a && b
			// a && (b || c)
			x, err := p.Parse(expr.X)
			if err != nil {
				return nil, err
			}
			y, err := p.Parse(expr.Y)
			if err != nil {
				return nil, err
			}
			return &LogicValidator{
				Qualifier: DefaultQualifier,
				Name:      "&&",
				Left:      x,
				Right:     y,
			}, nil

		case token.LOR:
			// a || b
			// a && b || c
			x, err := p.Parse(expr.X)
			if err != nil {
				return nil, err
			}
			y, err := p.Parse(expr.Y)
			if err != nil {
				return nil, err
			}
			return &LogicValidator{
				Qualifier: DefaultQualifier,
				Name:      "||",
				Left:      x,
				Right:     y,
			}, nil
		}

	case *ast.Ident:
		// a
		// _
		return &LeafValidator{
			Name: expr.Name,
		}, nil

	case *ast.CallExpr:
		switch fun := expr.Fun.(type) {
		case *ast.Ident:
			// a()
			var args []string
			for _, arg := range expr.Args {
				argValue, _, err := p.parseCallArgExpr(arg)
				if err != nil {
					return nil, err
				}
				args = append(args, argValue)
			}
			return &LeafValidator{
				Name: fun.Name,
				Args: args,
			}, nil

		case *ast.SelectorExpr:
			switch x := fun.X.(type) {
			case *ast.Ident:
				// a.b()
				msg, err := p.parseMsgExpr(expr)
				if err != nil {
					return nil, err
				}
				return &LeafValidator{
					Name: x.Name,
					Msg:  msg,
				}, nil

			case *ast.CallExpr:
				// a().b()
				ident, ok := x.Fun.(*ast.Ident)
				if !ok {
					return nil, p.error("", x)
				}

				msg, err := p.parseMsgExpr(expr)
				if err != nil {
					return nil, err
				}

				var args []string
				for _, arg := range x.Args {
					argValue, _, err := p.parseCallArgExpr(arg)
					if err != nil {
						return nil, err
					}
					args = append(args, argValue)
				}

				return &LeafValidator{
					Name: ident.Name,
					Args: args,
					Msg:  msg,
				}, nil

			default:
				return nil, p.error("", x)
			}

		default:
			return nil, p.error("", fun)
		}

	default:
		return nil, p.error("", expr)
	}

	return nil, nil
}

func (p Parser) parseCallArgExpr(e ast.Expr) (string, token.Token, error) {
	switch arg := e.(type) {
	case *ast.BasicLit:
		// gt(0)
		//    ^
		return arg.Value, arg.Kind, nil
	case *ast.Ident:
		// gt(min)
		//    ^^^
		return arg.Name, -1, nil
	default:
		return "", -1, p.error("", e)
	}
}

// parseMsgExpr extracts the custom error message from `msg("...")`.
func (p Parser) parseMsgExpr(e *ast.CallExpr) (string, error) {
	sel := e.Fun.(*ast.SelectorExpr)
	if sel.Sel.Name != "msg" {
		return "", p.error(p.string(sel.X)+".msg", sel)
	}

	if len(e.Args) != 1 {
		return "", p.error(p.string(sel.X)+".msg(\"...\")", sel)
	}

	msg, kind, err := p.parseCallArgExpr(e.Args[0])
	if err != nil {
		return "", err
	}

	if kind != token.STRING {
		return "", p.error("a string", e.Args[0])
	}

	return msg, nil
}

func (p Parser) string(e ast.Expr) string {
	start, end := e.Pos()-1, e.End()-1
	return p.S[start:end]
}

func (p Parser) error(expected string, e ast.Expr) error {
	if expected == "" {
		return fmt.Errorf("1:%d unexpected %s", e.Pos(), p.string(e))
	}
	return fmt.Errorf("1:%d expected %s, found %s", e.Pos(), expected, p.string(e))
}
