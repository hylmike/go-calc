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
	return fmt.Sprintf("NumberExpressionAST: %s", strconv.FormatFloat(n.Val, 'f', 0, 64))
}

func (b BinaryExprAST) toString() string {
	return fmt.Sprintf("BinaryExpressionAST: (%s %s %s", b.Operator, b.Lhs.toString(), b.Rhs.toString())
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

func (a *AST) geTokenPrecedence() int {
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
	a.getNextToken()
	return numberAST
}

func (a *AST) ParseExpression() ExpressionAST {
	lhs := a.parsePrimary()
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
			a.getNextToken()
			return expAST
		} else {
			return a.parseNumber()
		}
	default:
		return nil
	}
}

func (a *AST) parseBinaryOpRHS(minPrecedence int, lhs ExpressionAST) ExpressionAST {
	for {
		tokenPrecedence := a.geTokenPrecedence()
		if tokenPrecedence < minPrecedence {
			return lhs
		}

		binaryOp := a.currentToken.Tok
		if a.getNextToken() == nil {
			return lhs
		}
		rhs := a.parsePrimary()
		if rhs == nil {
			return nil
		}
		nextPrecedence := a.geTokenPrecedence()
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
