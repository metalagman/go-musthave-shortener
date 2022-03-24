/*
Package nomainexit allows you to check if you're making direct call to os.Exit
within the main function of main package of your application.
*/
package nomainexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// Analyzer instance
var Analyzer = &analysis.Analyzer{
	Name: "nomainexit",
	Doc:  "check for using os.Exit in function main of package main",
	Run:  run,
}

// run analyzer code
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

// isPkgDot checks provided selector matches the package name and function name
func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

// isIdent checks if provided expression is equal to provided ident
func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}
