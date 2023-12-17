// Package mathCalcEngine
// Read token parsed from a match expression
// Output AST and calculation result
package mathCalcEngine

import (
	"errors"
	"fmt"
	"strconv"
)

type ExpressionAST interface {
	toString() string
}

type NumberExprAST struct {
	Val float64
}

type BinaryExprAST struct {
	Operator string
	Lhs, Rhs ExpressionAST
}

func (n NumberExprAST) toString() string {
	return fmt.Sprintf("( %s )", strconv.FormatFloat(n.Val, 'f', 0, 64))
}

func (b BinaryExprAST) toString() string {
	return fmt.Sprintf("(L %s R)", b.Operator)
}

// math operators precedence, higher ones will be calculated first
var optPrecedence = map[string]int{"+": 10, "-": 10, "*": 20, "/": 20, "%": 20}

type AST struct {
	// results from lexical analysis
	Tokens []*Token
	// Original math expression
	source       string
	currentToken *Token
	currentIndex int
	Err          error
}

func (a *AST) getCurrentTokenPrecedence() int {
	if precedence, ok := optPrecedence[a.currentToken.Tok]; ok {
		return precedence
	}

	return -1
}

func (a *AST) getNextToken() *Token {
	a.currentIndex++
	if a.currentIndex < len(a.Tokens) {
		a.currentToken = a.Tokens[a.currentIndex]
		return a.currentToken
	}

	return nil
}

func (a *AST) parseNumber() NumberExprAST {
	number, err := strconv.ParseFloat(a.currentToken.Tok, 64)
	if err != nil {
		a.Err = errors.New(
			fmt.Sprintf(
				"%v\nShould be '(' or '0-9' but get '%s'\n%s",
				err.Error(),
				a.currentToken.Tok,
				ErrPos(a.source, a.currentToken.Offset),
			),
		)
		return NumberExprAST{}
	}

	numberAST := NumberExprAST{Val: number}
	return numberAST
}

func (a *AST) ParseExpression() ExpressionAST {
	lhs := a.parsePrimary()
	if a.getNextToken() == nil {
		return lhs
	}
	return a.parseBinaryOpRHS(0, lhs)
}

func (a *AST) parsePrimary() ExpressionAST {
	switch a.currentToken.Type {
	case Literal:
		return a.parseNumber()
	case Operator:
		// process () syntax
		if a.currentToken.Tok == "(" {
			a.getNextToken()
			expAST := a.ParseExpression()
			if expAST == nil {
				return nil
			}
			if a.currentToken.Tok != ")" {
				a.Err = errors.New(fmt.Sprintf(
					"Should be ')' but get %s\n%s",
					a.currentToken.Tok,
					ErrPos(a.source, a.currentToken.Offset),
				))
				return nil
			}
			return expAST
		} else {
			// Means this position should be '(' or number, use parNumber() to throw error
			return a.parseNumber()
		}
	default:
		return nil
	}
}

func (a *AST) parseBinaryOpRHS(minPrecedence int, lhs ExpressionAST) ExpressionAST {
	for {
		tokenPrecedence := a.getCurrentTokenPrecedence()
		if tokenPrecedence < minPrecedence {
			return lhs
		}

		binaryOp := a.currentToken.Tok
		if a.getNextToken() == nil {
			return lhs
		}
		rhs := a.parsePrimary()
		if a.Err != nil || rhs == nil {
			return nil
		}
		if a.getNextToken() == nil {
			return BinaryExprAST{
				Operator: binaryOp,
				Lhs:      lhs,
				Rhs:      rhs,
			}
		}
		nextPrecedence := a.getCurrentTokenPrecedence()
		if tokenPrecedence < nextPrecedence {
			rhs = a.parseBinaryOpRHS(tokenPrecedence+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		lhs = BinaryExprAST{
			Operator: binaryOp,
			Lhs:      lhs,
			Rhs:      rhs,
		}
	}
}

func CreateAST(tokens []*Token, source string) *AST {
	ast := &AST{
		Tokens: tokens,
		source: source,
	}

	if ast.Tokens == nil || len(ast.Tokens) == 0 {
		ast.Err = errors.New("empty token list")
	} else {
		ast.currentIndex = 0
		ast.currentToken = ast.Tokens[0]
	}

	return ast
}

func GetMaxLevel(root ExpressionAST, level int) int {
	if root == nil {
		return level
	}

	switch root.(type) {
	case NumberExprAST:
		return level + 1
	case BinaryExprAST:
		return max(GetMaxLevel(root.(BinaryExprAST).Lhs, level+1), GetMaxLevel(root.(BinaryExprAST).Rhs, level+1))
	default:
		return level
	}
}

func PrintCalcAST(root ExpressionAST) {
	if root == nil {
		return
	}

	queue := []ExpressionAST{root}

	for len(queue) > 0 {
		layerSize := len(queue)

		for i := 0; i < layerSize; i++ {
			node := queue[0]
			queue = queue[1:]

			switch node.(type) {
			case BinaryExprAST:
				lNode := node.(BinaryExprAST).Lhs
				rNode := node.(BinaryExprAST).Rhs
				if lNode != nil {
					queue = append(queue, lNode)
				}
				if rNode != nil {
					queue = append(queue, rNode)
				}
			case NumberExprAST:
			}
			fmt.Printf("%s\t", node.toString())
		}
		fmt.Println()
	}
}
