package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/inspector"
)

func main() {
	// ソースコードを文字列で定義
	const src = `package main
	func main() {
		println("hello")
	}`

	// トークン位置情報を保持する構造体
	fset := token.NewFileSet()

	// ソースコードを構文解析して*ast.Fileを作る
	f, err := parser.ParseFile(fset, "my.go", src, 0)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	// inspectorを作成
	inspect := inspector.New([]*ast.File{f})

	// 探索したいノードの型を指定(ここでは関数呼び出し)
	typs := []ast.Node{new(ast.CallExpr)}

	inspect.Preorder(typs, func(n ast.Node) {
		// 型を表示
		fmt.Printf("見つけたノードの型: %T\n", n)

		// n.(*ast.CallExpr)に変換して関数名などを表示
		call := n.(*ast.CallExpr)
		if funIdent, ok := call.Fun.(*ast.Ident); ok {
			fmt.Printf("関数名: %s\n", funIdent.Name)
		}
	})
}

