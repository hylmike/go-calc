package mathCalcEngine

import (
	"errors"
	"fmt"
)

func GetExprASTResult(root ExpressionAST) float64 {
	var lVal, rVal float64
	switch root.(type) {
	case BinaryExprAST:
		ast := root.(BinaryExprAST)
		lVal = GetExprASTResult(ast.Lhs)
		rVal = GetExprASTResult(ast.Rhs)

		switch ast.Operator {
		case "+":
			return lVal + rVal
		case "-":
			return lVal - rVal
		case "*":
			return lVal * rVal
		case "/":
			if rVal == 0 {
				panic(errors.New(fmt.Sprintf("division by zero: '%f/%f'", lVal, rVal)))
			}
			return lVal / rVal
		case "%":
			return float64(int(lVal) / int(rVal))
		default:
			panic(errors.New(fmt.Sprintf("invalid operator: '%s'", ast.Operator)))
		}
	case NumberExprAST:
		return root.(NumberExprAST).Val
	}

	return 0.0
}
