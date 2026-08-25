package document

// import (
// 	"slices"
// 	"testing"
// )

// func TestDocumentSlice(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		src  []byte
// 		rng  Range
// 		want []byte
// 	}{
// 		{
// 			name: "First Test",
// 			src:  []byte{'h', 'o', 'l', 'a'},
// 			rng:  NewRange(0, 1),
// 			want: []byte{'h'},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			dc, err := New(tt.src)
// 			if err != nil {
// 				t.Fatalf("document.New(%v) = %v; want nil", tt.src, err)
// 			}

// 			got, err := dc.Slice(tt.rng)
// 			if err != nil {
// 				t.Fatalf("document.Slice(%v) = %v; want nil", tt.rng, err)
// 			}

// 			if !slices.Equal(got, tt.want) {
// 				t.Errorf(
// 					"document.Slice(%v) = %v; want = %v",
// 					tt.rng,
// 					got,
// 					tt.want,
// 				)
// 			}
// 		})
// 	}
// }
