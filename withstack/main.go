package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/inspector"
)

func main() {
	// コード対象
	const src = `
		package main

		import "fmt"

		func main() {
			println("hello")
			fmt.Println("world")
			x := add(1, 2)
		}

		func add(a, b int) int {
			return a + b
		}
	`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example.go", src, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	// inspectorを作成
	ins := inspector.New([]*ast.File{file})

	// 関数呼び出しノード(CallExpr)を対象に走査
	typs := []ast.Node{new(ast.CallExpr)}
	ins.WithStack(typs, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}

		fmt.Println("CallExprの親ノード")
		for _, node := range stack {
			fmt.Printf(" %T\n", node)
		}

		// CallExprの情報を出力
		call := n.(*ast.CallExpr)
		fmt.Printf("呼び出し先: ")
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			fmt.Printf("%s\n", fun.Name)
		case *ast.SelectorExpr:
			if pkg, ok := fun.X.(*ast.Ident); ok {
				fmt.Printf("%s.%s\n", pkg.Name, fun.Sel.Name)
			} else {
				fmt.Printf("%#v\n", fun)
			}
		default:
			fmt.Printf("%#v\n", fun)
		}
		fmt.Println()
		return true
	})
}