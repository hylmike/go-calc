// Package mathCalcEngine
// Read math expression like '3+5*(3+4)-9/3'
// Output token list
package mathCalcEngine

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Literal = iota
	Operator
)

type Token struct {
	// Original token string
	Tok string
	// Literal, Operator
	Type int
	// tok position, used for print error position
	Offset int
}

// Parser - Lexical analyzer struct
type Parser struct {
	// original math expression string
	Source string
	// current scanning char
	Char byte
	// scan location
	offset int
	// error from string parsing
	Err error
}

func (p *Parser) nextCh() error {
	p.offset++
	if p.offset < len(p.Source) {
		p.Char = p.Source[p.offset]
		return nil
	}

	// reach the end of string
	return errors.New("EOF")
}

func (p *Parser) isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\v' || ch == '\f' || ch == '\r'
}

func (p *Parser) isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.' || ch == '_' || ch == 'e'
}

func ErrPos(s string, pos int) string {
	r := strings.Repeat("-", len(s)) + "\n"
	s += "\n"
	for i := 0; i < pos; i++ {
		s += " "
	}
	s += "^\n"

	return r + s + r
}

func (p *Parser) nextTok() *Token {
	if p.offset >= len(p.Source) || p.Err != nil {
		return nil
	}

	var err error
	for p.isWhiteSpace(p.Char) && err == nil {
		err = p.nextCh()
	}

	start := p.offset
	var token *Token
	switch p.Char {
	case '(', ')', '+', '-', '*', '/', '%':
		token = &Token{
			Tok:    string(p.Char),
			Type:   Operator,
			Offset: start,
		}
		err = p.nextCh()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Here implicitly includes p.nextCh()
		for p.isNumber(p.Char) && p.nextCh() == nil {
		}
		token = &Token{
			Tok:    strings.ReplaceAll(p.Source[start:p.offset], "_", ""),
			Type:   Literal,
			Offset: start,
		}
	default:
		message := fmt.Sprintf(
			"Symbol error: unknown '%v', position [%v:]\n%s",
			string(p.Char),
			start,
			ErrPos(p.Source, start),
		)
		p.Err = errors.New(message)
	}

	return token
}

func (p *Parser) parse() []*Token {
	tokens := make([]*Token, 0)

	for {
		token := p.nextTok()
		if token == nil {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func Parse(source string) ([]*Token, error) {
	p := &Parser{
		Source: source,
		Err:    nil,
		Char:   source[0],
	}

	tokens := p.parse()
	if p.Err != nil {
		return nil, p.Err
	}

	return tokens, nil
}
