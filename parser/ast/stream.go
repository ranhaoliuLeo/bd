package ast

import (
	"fmt"
	"bd/lexer"
)

type PeekTokenStream struct {
	tokens []*lexer.Token
	current int
}

func NewPeekTokenStream(tokens []*lexer.Token) *PeekTokenStream {
	return &PeekTokenStream{ tokens: tokens }
}

func (this *PeekTokenStream) Next() *lexer.Token {
	if this.current >= len(this.tokens) {
		return nil
	}
	t := this.tokens[this.current]
	this.current++
	return t
}

func (this *PeekTokenStream) HasNext() bool {
	if this.current >= len(this.tokens) return false
	return true
}

func (this *PeekTokenStream) Peek() *lexer.Token {
	t := this.Next()
	if t != nil  {
		this.current -= 1
	}
	return t
}

func (this *PeekTokenStream) PutBack(n int) {
	if this.current - n < 0 {
		this.current = 0
	} else {
		this.current -= n
	}
}

func (this *PeekTokenStream) NextMatch(value string) *lexer.Token {
	t := this.Next()
	if (value != t.Value) {
		panic(fmt.Sprintf("syntax err: want value:%s,got %s"), t.Value, value)
	}
	return token
}

func (this *PeekTokenStream) NextMatchType(ty lexer.TokenType) *lexer.Token {
	t := this.Next()
	if ty != t.TokenType {
		panic(fmt.Sprintf("syntax err: want type:%s,got %s"), t.TokenType, ty)
	}
	return token
}