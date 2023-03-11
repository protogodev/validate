package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

const (
	DefaultQualifier = "validating"
)

type Validator interface {
	// SetQualifier sets the validator's qualifier to q.
	SetQualifier(q string)

	// ConvertName convert the validator's name (recursively) by applying
	// the corresponding converter to each validator.
	ConvertName(converters map[string]Converter, typ types.Type) error

	// ExprString returns the validating-style expression string.
	ExprString() string
}

// LeafValidator is an expression that represents a leaf validator.
type LeafValidator struct {
	Qualifier string
	Name      string
	Type      string
	Args      []string
	Msg       string
}

func (v *LeafValidator) SetQualifier(q string) {
	v.Qualifier = q
}

func (v *LeafValidator) ConvertName(converters map[string]Converter, typ types.Type) error {
	convert, ok := converters[v.Name]
	if !ok {
		return fmt.Errorf("found no converter for %q", v.Name)
	}

	name, err := convert(v.Name, typ)
	if err != nil {
		return err
	}

	v.Name = name
	return nil
}

func (v *LeafValidator) ExprString() string {
	args := strings.Join(v.Args, ", ")
	if strings.ToLower(v.Name) == "match" {
		args = fmt.Sprintf("regexp.MustCompile(%s)", args)
	}

	if v.Msg == "" {
		return fmt.Sprintf("%s.%s(%s)", v.Qualifier, v.Name, args)
	}
	return fmt.Sprintf("%s.%s(%s).Msg(%s)", v.Qualifier, v.Name, args, v.Msg)
}

// LogicValidator is an expression that represents a logic validator (i.e. `Not`, `And/All` or `Or/Any`).
type LogicValidator struct {
	Qualifier string
	Name      string
	Left      Validator
	Right     Validator // `Not` has no Right validator.
}

func (v *LogicValidator) SetQualifier(q string) {
	v.Qualifier = q
	v.Left.SetQualifier(q)
	if v.Right != nil {
		v.Right.SetQualifier(q)
	}
}

func (v *LogicValidator) ConvertName(converters map[string]Converter, typ types.Type) error {
	switch v.Name {
	case "!":
		v.Name = "Not"
	case "&&":
		v.Name = "All"
	case "||":
		v.Name = "Any"
	}

	if err := v.Left.ConvertName(converters, typ); err != nil {
		return err
	}

	if v.Right != nil {
		return v.Right.ConvertName(converters, typ)
	}

	return nil
}

func (v *LogicValidator) ExprString() string {
	if v.Right != nil {
		return fmt.Sprintf("%s.%s(%s, %s)", v.Qualifier, v.Name, v.Left.ExprString(), v.Right.ExprString())
	}
	return fmt.Sprintf("%s.%s(%s)", v.Qualifier, v.Name, v.Left.ExprString())
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
				Name:      "!",
				Qualifier: DefaultQualifier,
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
			Qualifier: DefaultQualifier,
			Name:      expr.Name,
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
				Qualifier: DefaultQualifier,
				Name:      fun.Name,
				Args:      args,
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
					Qualifier: DefaultQualifier,
					Name:      x.Name,
					Msg:       msg,
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
					Qualifier: DefaultQualifier,
					Name:      ident.Name,
					Args:      args,
					Msg:       msg,
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
