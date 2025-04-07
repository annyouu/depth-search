// 抽象構文木を深さ優先探索で走査する方法

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()

	src := `
		package main

		import "fmt"

		// Add は 2 つの整数を足し合わせる関数
		func Add(x, y int) int {
    		sum := x + y
    		fmt.Println("sum:", sum)
    		return sum
		}
	`

	fileNode, err := parser.ParseFile(fset, "example.go", src, parser.ParseComments)
	if err != nil {
		fmt.Println("パースエラー:", err)
		return
	}

	// ast.InspectでASTを深さ優先探索する
	ast.Inspect(fileNode, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		// ノードを型ごとに処理を分岐する
		switch node := n.(type) {
		case *ast.FuncDecl:
			fmt.Printf("関数宣言: name=%s, params=%d, results=%d\n",
				node.Name.Name,
				len(node.Type.Params.List),
				func() int {
					if node.Type.Results != nil {
						return len(node.Type.Results.List)
					}
					return 0
				}(),
			)
		case *ast.AssignStmt:
			// 代入文ノードを発見
			fmt.Printf("代入文: %s\n", node.Tok) // Tokは := や = を表す
		case *ast.CallExpr:
			// 関数呼び出しノードを発見
			if fun, ok := node.Fun.(*ast.SelectorExpr); ok {
				fmt.Printf("メソッド呼び出し: %s.%s\n", fun.X, fun.Sel.Name)
			} 
		}
		return true
	})

}
