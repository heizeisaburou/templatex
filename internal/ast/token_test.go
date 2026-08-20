package ast

import (
	"errors"
	"testing"

	"github.com/heizeisaburou/templatex/internal/document"
)

func TestNewTokenConstructor(t *testing.T) {
	tests := []struct {
		name string
		kind TokenKind
		rng  document.Range
		want Token
	}{
		{
			name: "unknown",
			kind: TokenUnknown,
			rng:  document.NewRange(0, 1),
			want: Token{TokenUnknown, document.NewRange(0, 1)},
		},
		{
			name: "whitespace_standard_range",
			kind: TokenWhitespace,
			rng:  document.NewRange(1, 4),
			want: Token{TokenWhitespace, document.NewRange(1, 4)},
		},
		{
			name: "whitespace_large_range",
			kind: TokenWhitespace,
			rng:  document.NewRange(100, 250),
			want: Token{TokenWhitespace, document.NewRange(100, 250)},
		},
		{
			name: "limit_boundary_range",
			kind: TokenLimit,
			rng:  document.NewRange(5, 5),
			want: Token{TokenLimit, document.NewRange(5, 5)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewToken(tt.kind, tt.rng)

			if err != nil {
				t.Fatalf(
					"\nNewToken(%v, %v) = %v; want=%v, err=%v",
					tt.kind,
					tt.rng,
					got,
					tt.want,
					err,
				)
			}

			if got != tt.want {
				t.Errorf(
					"\nNewToken(%v, %v) = %v; want=%v",
					tt.kind,
					tt.rng,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestNewTokenConstructorBadPath(t *testing.T) {
	tests := []struct {
		name string
		kind TokenKind
		rng  document.Range
		err  error
	}{
		{
			name: "unknown",
			kind: -1,
			rng:  document.NewRange(0, 1),
			err:  ErrNotValidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewToken(tt.kind, tt.rng)

			if !errors.Is(err, tt.err) {
				t.Errorf(
					"\nNewToken(%v, %v) = %v; err=%v",
					tt.kind,
					tt.rng,
					got,
					err,
				)
			}
		})
	}
}
