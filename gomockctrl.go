package gomockctrl

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "gomockctrl",
	Doc:  "gomockctrl is a linter that detects gomock.Controller defined outside t.Run used inside it",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	var (
		lastController    *ast.CallExpr
		lastControllerEnd token.Pos
		stack             []token.Pos
	)

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		for len(stack) > 0 && n.Pos() >= stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		}

		if n.Pos() >= lastControllerEnd {
			lastController = nil
		}

		switch n := n.(type) {
		case *ast.CallExpr:
			fun, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			pkgName, ok := fun.X.(*ast.Ident)
			if !ok {
				return
			}
			switch pkgName.Name + "." + fun.Sel.Name {
			case "gomock.NewController":
				lastController = n
				lastControllerEnd = stack[len(stack)-1]
			case "t.Run":
				if lastController != nil && lastController.Pos() < n.Pos() {
					pass.Report(analysis.Diagnostic{
						Pos:     n.Pos(),
						Message: "testing.T.Run was run after gomock.NewController inside the same test function",
					})
				}
			}
		case *ast.FuncDecl, *ast.FuncLit:
			stack = append(stack, n.End())
		}
	})
	return nil, nil
}
