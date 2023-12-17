package main

import (
	"calc/mathCalcEngine"
	"fmt"
	"os"
)

func main() {
	exp := os.Args[1]

	tokens, err := mathCalcEngine.Parse(exp)
	if err != nil {
		fmt.Println("Invalid expression: ", err.Error())
		return
	}

	ast := mathCalcEngine.CreateAST(tokens, exp)
	if ast.Err != nil {
		fmt.Println("Invalid expression: ", ast.Err.Error())
		return
	}

	astTree := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("Invalid expression: ", ast.Err.Error())
		return
	}

	fmt.Printf("Expression AST: %+v\n", astTree)

	result := mathCalcEngine.GetExprASTResult(astTree)

	fmt.Printf("The calculation result of math expression '%s' is %.2f\n", exp, result)
}
