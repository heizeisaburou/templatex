package ast

import (
	"testing"

	"github.com/heizeisaburou/templatex/internal/document"
)

func newToken(t *testing.T, kind TokenKind, rng document.Range) Token {
	t.Helper()
	token, err := NewToken(kind, rng)
	if err != nil {
		t.Fatal("newToken() error al crear token")
	}

	return token
}

func newRange(t *testing.T, start document.ByteOffset, end document.ByteOffset) document.Range {

	t.Helper()
	rng, err := document.NewRange(start, end)
	if err != nil {
		t.Fatal("newRange() rango invalido")
	}

	return rng
}

func TestTokenFormatterString(t *testing.T) {
	tests := []struct {
		name  string
		src   []byte
		kind  TokenKind
		start document.ByteOffset
		end   document.ByteOffset
		want  string
	}{
		{
			name:  "normal case (standard identifier)",
			src:   []byte("let x = 10;"),
			kind:  1, // In your codebase, 1 resolves to TokenWhitespace
			start: 0,
			end:   3,
			want:  `Token{kind=TokenWhitespace, rng=[0, 3), src="let"}`,
		},
		{
			name:  "edge case: trailing space and newline",
			src:   []byte("func \n"),
			kind:  TokenUnknown, // Using a valid kind (0) so NewToken doesn't return an error
			start: 4,
			end:   6,
			want:  `Token{kind=TokenUnknown(0), rng=[4, 6), src=" \n"}`,
		},
		{
			name:  "edge case: internal quotes",
			src:   []byte(`"hello"`),
			kind:  TokenUnknown, // Using a valid kind (0) so NewToken doesn't return an error
			start: 0,
			end:   7,
			want:  `Token{kind=TokenUnknown(0), rng=[0, 7), src="\"hello\""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := document.New(tt.src)

			if err != nil {
				t.Fatalf("document.New(%v) unexpected_error: %v", tt.src, err)
			}

			rng := newRange(t, tt.start, tt.end)
			token := newToken(t, tt.kind, rng)

			formatter := NewTokenFormatter(doc)

			got, err := formatter.Sprint(token)

			if err != nil {
				t.Fatalf(
					"TokenFormatter.Sprint(%v) unexpected_error: %v",
					token,
					err,
				)
			}

			if got != tt.want {
				t.Errorf(
					"TokenFormatter.Sprint(%v) = \"%s\"; want = \"%s\"",
					token,
					got,
					tt.want,
				)
			}

		})
	}
}
