package ast

import (
	"errors"

	"github.com/heizeisaburou/templatex/internal/document"
)

type TokenKind int

var (
	ErrNotValidToken = errors.New("El token utilizado no es valido")
)

const (
	TokenUnknown TokenKind = iota
	TokenWhitespace
	TokenLimit
)

type Token struct {
	kind TokenKind
	rng  document.Range
}

func NewToken(kind TokenKind, rng document.Range) (Token, error) {

	if kind < TokenUnknown || kind > TokenLimit {
		return Token{}, ErrNotValidToken
	}

	return Token{kind: kind, rng: rng}, nil
}

func (t Token) Range() document.Range {
	return t.rng
}

func (t Token) Kind() TokenKind {
	return t.kind
}
