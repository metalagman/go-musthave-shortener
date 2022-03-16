package nomainexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nomainexit",
	Doc:  "check for using os.Exit in function main of package main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isPkgDot(call.Fun, "os", "Exit") {
				pass.Reportf(x.Pos(), "os.Exit called directly in func main package main")
			}
		}
	}

	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				return x.Name.Name == "main"
			case *ast.FuncDecl:
				return x.Name.Name == "main"
			case *ast.ExprStmt: // выражение
				expr(x)
			}
			return true
		})
	}
	return nil, nil
}

// helpers
// =======
func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}
