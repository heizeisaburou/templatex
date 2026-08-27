package ast

// import (
// 	"errors"
// 	"testing"

	"github.com/heizeisaburou/templatex/internal/document"
)

type Range = document.Range

func newToken(t *testing.T, kind TokenKind, rng Range) Token {
  t.Helper()

  token, err := NewToken(kind, rng)
  if err != nil {
    t.Fatalf("newToken() error al crear token")
  }

  return token
}

// func TestNewTokenConstructor(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		kind TokenKind
// 		rng  document.Range
// 		want Token
// 	}{
// 		{
// 			name: "unknown",
// 			kind: TokenUnknown,
// 			rng:  document.NewRange(0, 1),
// 			want: Token{TokenUnknown, document.NewRange(0, 1)},
// 		},
// 		{
// 			name: "whitespace_standard_range",
// 			kind: TokenWhitespace,
// 			rng:  document.NewRange(1, 4),
// 			want: Token{TokenWhitespace, document.NewRange(1, 4)},
// 		},
// 		{
// 			name: "whitespace_large_range",
// 			kind: TokenWhitespace,
// 			rng:  document.NewRange(100, 250),
// 			want: Token{TokenWhitespace, document.NewRange(100, 250)},
// 		},
// 		{
// 			name: "limit_boundary_range",
// 			kind: tokenLimit,
// 			rng:  document.NewRange(5, 5),
// 			want: Token{tokenLimit, document.NewRange(5, 5)},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := NewToken(tt.kind, tt.rng)

// 			if err != nil {
// 				t.Fatalf(
// 					"\nNewToken(%v, %v) = %v; want=%v, err=%v",
// 					tt.kind,
// 					tt.rng,
// 					got,
// 					tt.want,
// 					err,
// 				)
// 			}

// 			if got != tt.want {
// 				t.Errorf(
// 					"\nNewToken(%v, %v) = %v; want=%v",
// 					tt.kind,
// 					tt.rng,
// 					got,
// 					tt.want,
// 				)
// 			}
// 		})
// 	}
// }

// func TestNewTokenConstructorBadPath(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		kind TokenKind
// 		rng  document.Range
// 		err  error
// 	}{
// 		{
// 			// Edge Case: The immediate lower bound (already in your draft, kept for completeness)
// 			name: "just_below_unknown",
// 			kind: -1,
// 			rng:  document.NewRange(0, 1),
// 			err:  ErrNotValidToken,
// 		},
// 		{
// 			// Extreme Case: Deeply negative value
// 			name: "far_below_unknown",
// 			kind: -999,
// 			rng:  document.NewRange(0, 1),
// 			err:  ErrNotValidToken,
// 		},
// 		{
// 			// Edge Case: The immediate upper bound (TokenLimit + 1)
// 			// TokenLimit is 2, so this tests a kind of 3.
// 			name: "just_above_limit",
// 			kind: tokenLimit + 1,
// 			rng:  document.NewRange(0, 1),
// 			err:  ErrNotValidToken,
// 		},
// 		{
// 			// Extreme Case: High positive value beyond the limit
// 			name: "far_above_limit",
// 			kind: 999,
// 			rng:  document.NewRange(10, 20),
// 			err:  ErrNotValidToken,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := NewToken(tt.kind, tt.rng)

// 			if !errors.Is(err, tt.err) {
// 				t.Errorf(
// 					"\nNewToken(%v, %v) = %v; err=%v",
// 					tt.kind,
// 					tt.rng,
// 					got,
// 					err,
// 				)
// 			}
// 		})
// 	}
// }
