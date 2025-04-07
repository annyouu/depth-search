package main

import (
	"fmt"
	"go/ast"
	"go/parser"
)

// ast.Vistorinterfaceを関数で実装するための型
type VisitFunc func(n ast.Node) ast.Visitor

// Visitメソッドを定義することでVisitFuncがast.Visitorinterfaceを満たす
func (v VisitFunc) Visit(n ast.Node) ast.Visitor {
	return v(n)
}

func main() {
	// 解析対象の式をパース
	expr, err := parser.ParseExpr("v+1")
	if err != nil {
		fmt.Println("パースエラー:", err)
		return
	}

	// Visitorを定義（二項演算式を見つけたら、その内部の識別子を探す別のVisitorに移行する）
	var visitor ast.Visitor
	visitor = VisitFunc(func(n ast.Node) ast.Visitor {
		// ノードが*ast.BinaryExprでなければ、同じvisitorを継続する
		if _, ok := n.(*ast.BinaryExpr); !ok {
			return visitor
		}

		// BinaryExprを見つけた → 次は識別子を探すvisitorに切り替える
		return VisitFunc(func(n ast.Node) ast.Visitor {
			// ノードが識別子(Ident)なら名前を出力
			if ident, ok := n.(*ast.Ident); ok {
				fmt.Println("見つかった識別子:", ident.Name)
			}
			// さらに子ノードも探索し続ける
			return nil
		})
	})

	ast.Walk(visitor, expr)
}