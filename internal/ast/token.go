package ast

import "github.com/heizeisaburou/templatex/internal/document"

type TokenKind int

const (
	TokenWhitespace TokenKind = iota
)

type Token struct {
	kind TokenKind
	rng  document.Range
}

func NewToken(kind TokenKind, rng document.Range) Token {
}

func (t Token) Range() document.Range {
}

func (t Token) Kind() TokenKind {
}
