package ast

import (
	"errors"
	"testing"

	"github.com/heizeisaburou/templatex/internal/document"
)

func TestNewTokenConstructor(t *testing.T) {
	tests := []struct {
		name      string
		kind      TokenKind
		start     document.ByteOffset
		end       document.ByteOffset
		wantKind  TokenKind
		wantStart document.ByteOffset
		wantEnd   document.ByteOffset
	}{
		{
			name:      "unknown",
			kind:      TokenUnknown,
			start:     0,
			end:       1,
			wantKind:  TokenUnknown,
			wantStart: 0,
			wantEnd:   1,
		},
		{
			name:      "whitespace_standard_range",
			kind:      TokenWhitespace,
			start:     1,
			end:       4,
			wantKind:  TokenWhitespace,
			wantStart: 1,
			wantEnd:   4,
		},
		{
			name:      "whitespace_large_range",
			kind:      TokenWhitespace,
			start:     100,
			end:       250,
			wantKind:  TokenWhitespace,
			wantStart: 100,
			wantEnd:   250,
		},
		{
			name:      "limit_boundary_range",
			kind:      tokenLimit,
			start:     5,
			end:       5,
			wantKind:  tokenLimit,
			wantStart: 5,
			wantEnd:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := newRange(t, tt.start, tt.end)
			wantRng := newRange(t, tt.wantStart, tt.wantEnd)

			gotToken := newToken(t, tt.kind, rng)
			wantToken := newToken(t, tt.wantKind, wantRng)

			if gotToken != wantToken {
				t.Errorf(
					"\nNewToken(%v, %v) = %v; want=%v",
					tt.kind,
					rng,
					gotToken,
					wantToken,
				)
			}
		})
	}
}

func TestNewTokenConstructorBadPath(t *testing.T) {
	tests := []struct {
		name  string
		kind  TokenKind
		start document.ByteOffset
		end   document.ByteOffset
		err   error
	}{
		{
			// Edge Case: el límite inferior inmediato.
			name:  "just_below_unknown",
			kind:  -1,
			start: 0,
			end:   1,
			err:   ErrNotValidToken,
		},
		{
			// Extreme Case: valor negativo muy alejado.
			name:  "far_below_unknown",
			kind:  -999,
			start: 0,
			end:   1,
			err:   ErrNotValidToken,
		},
		{
			// Edge Case: el límite superior inmediato (tokenLimit + 1).
			name:  "just_above_limit",
			kind:  tokenLimit + 1,
			start: 0,
			end:   1,
			err:   ErrNotValidToken,
		},
		{
			// Extreme Case: valor positivo muy por encima del límite.
			name:  "far_above_limit",
			kind:  999,
			start: 10,
			end:   20,
			err:   ErrNotValidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := newRange(t, tt.start, tt.end)

			got, err := NewToken(tt.kind, rng)

			if !errors.Is(err, tt.err) {
				t.Errorf(
					"\nNewToken(%v, %v) err = %v; want = %v",
					tt.kind,
					rng,
					err,
					tt.err,
				)
			}

			if got != (Token{}) {
				t.Errorf(
					"\nNewToken(%v, %v) = %v; want = %v",
					tt.kind,
					rng,
					got,
					Token{},
				)
			}
		})
	}
}
