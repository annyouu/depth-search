package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/inspector"
)

func main() {
	const src = "package main\n func main() { println(`hello`) }"

	fset := token.NewFileSet()

	// Goのコードをパースして抽象構文木(AST)を生成
	f, err := parser.ParseFile(fset, "my.go", src, 0)
	if err != nil {
		panic(err)
	}

	// inspectorの作成
	inspect := inspector.New([]*ast.File{f})

	// CallExprのみを走査
	typs := []ast.Node{new(ast.CallExpr)}

	// ノードに辿りつくまでの親ノードのスタックも取得する
	inspect.WithStack(typs, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}

		fmt.Println("親ノードのスタック:")
		for i, node := range stack {
			fmt.Printf("[%d] %T\n", i, node)
		}

		// ノードの位置情報と型を表示
		pos := fset.Position(n.Pos())
		fmt.Printf("この関数呼び出しは%sにあります\n", pos)
		fmt.Println("---")
		return true
	})
}