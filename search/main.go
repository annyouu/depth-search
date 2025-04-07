package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// 解析するソース
	const src = `
		package main

		// 通常の関数
		func Gopher() {
		}

		// レシーバ付きのメソッド（除外対象）
		func (r *Receiver) Gopher() {
		}

		// 変数（除外対象）
		var Gopher = 42
	`

	// ソースをパースしてASTを得る
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "example.go", src, 0)
	if err != nil {
		panic(err)
	}

	// ast.Inspectで全ノードを深さ優先探索
	ast.Inspect(file, func(n ast.Node) bool {
		// ノードが*ast.Identじゃないからスキップ
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		// 名前がGopherでないならスキップする
		if ident.Name != "Gopher" {
			return true
		}

		// 宣言元がないからスキップする
		if ident.Obj == nil {
			return true
		}

		// obj.Kindがast.Fun(関数宣言)であるかチェックする
		if ident.Obj.Kind == ast.Fun {
			// Obj.Declが*ast.FuncDeclかつ
			if fn, ok := ident.Obj.Decl.(*ast.FuncDecl); ok && fn.Recv == nil {
				pos := fset.Position(ident.Pos())
				fmt.Printf("関数 Gopherを発見: %s\n", pos)
			}
		}
		return true
	})
}