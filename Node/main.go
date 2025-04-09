package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/inspector"
)

func main() {
	// 解析対象のコード
	const src = `
		package main

		import "fmt"

		func main() {
			fmt.Println("Hello, world!")
			x := add(1, 2)
			y := multiply(x, 3)
			println(y)
		}

		func add(a, b int) int {
			return a + b
		}

		func multiply(a, b int) int {
			return a * b
		}
	`

	// ファイルセットを作成
	fset := token.NewFileSet()

	// ソースコードをパースしてASTを取得
	file, err := parser.ParseFile(fset, "example.go", src, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	// inspecotorを作成
	ins := inspector.New([]*ast.File{file})

	// 関数呼び出しノード(ast.CallExpr)を探す
	typs := []ast.Node{new(ast.CallExpr)}
	ins.Nodes(typs, func(n ast.Node, push bool) bool {
		if push {
			call := n.(*ast.CallExpr)
			if funIdent, ok := call.Fun.(*ast.Ident); ok {
				pos := fset.Position(call.Lparen)
				fmt.Printf("関数呼び出し: %s at %s\n", funIdent.Name, pos)
			} else {
				pos := fset.Position(call.Lparen)
				fmt.Printf("関数呼び出し(複雑な式) at %s: %#v\n", pos, call.Fun)
			}
		}
		return true
	})
}
