package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	// 解析対象の式
	expr, err := parser.ParseExpr(`a + b`)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	// ASTノードをApplyを使って走査&置換する
	n := astutil.Apply(expr, func(cr *astutil.Cursor) bool {
		switch cr.Name() {
		case "X":
			// 左辺を10に置き換え
			cr.Replace(&ast.BasicLit{
				Kind: token.INT,
				Value: "10",
			})
		case "Y":
			// 右辺を整数リテラル20に置き換える
			cr.Replace(&ast.BasicLit{
				Kind: token.INT,
				Value: "20",
			})
		}
		return true
	}, nil)

	// 置換後のASTを出力
	fmt.Println("=== AST Print ===")
	ast.Print(nil, n)
}