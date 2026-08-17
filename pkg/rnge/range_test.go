package rnge

import (
	"testing"
)

type Int int

func (i Int) Compare(other Int) int {
	switch {
	case i < other:
		return -1
	case i > other:
		return 1
	default:
		return 0
	}
}

func TestNew(t *testing.T) {
	t.Run("start menor que end", func(t *testing.T) {
		r := New[Int](5, 10)
		if got := r.Start(); got != 5 {
			t.Errorf("Start() = %d, want 5", got)
		}
		if got := r.End(); got != 10 {
			t.Errorf("End() = %d, want 10", got)
		}
	})

	t.Run("start mayor que end", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("New() must panic when start is bigger than end")
			}
		}()

		New[Int](10, 5)
	})
}

func TestSet(t *testing.T){
  t.Run("Set modifica los valores del rango y los índices los establece correctamente", func(t *testing.T){
    r := New[Int](2, 4)
    r.Set(5, 10)

    start := r.Start()
    end := r.End()

    if got := start; got != 5 {
    t.Errorf("start = %d, want 5", got)
    }
		if got := end; got != 10 {
			t.Errorf("end = %d, want 10", got)
		}
  })

  t.Run("Set entra en pánico si end es menor que start", func(t *testing.T){
		defer func() {
			if recover() == nil {
				t.Error("New() must panic when start is bigger than end")
			}
		}()

    wr := New[Int](2, 4)
    wr.Set(10, 5)
  })
}

func TestSetRange(t *testing.T){
  t.Run("SetRange si asigna a r los índices de other", func(t *testing.T){
    r := New[Int](2, 4)
    r.Set(5, 10)
    
    other := New[Int](6, 12)
    other.Set(10, 20)

    r.SetRange(other)

    if got := r.Start(); got != other.Start(){
      t.Errorf("start = %d, want 10", got)
    }

    if got := r.End(); got != other.End(){
      t.Errorf("end = %d, want 20", got)
    }
  })
}

func TestEmpty(t *testing.T){
  t.Run("Empty detecta si el rango está con valores nulos definidos", func(t* testing.T){
    r := New[Int](5, 10)
    r.Set(0, 0)

    e := r.start.Compare(r.end)

    if got := e; got != 0{
      t.Errorf("Range is empty (%b), want start and end different to zero", got)
    }
  })

  t.Run("", func(t *testing.T){
    var r Int
    if got := r; got != 0{
      t.Errorf("Range is empty (%b), want start and end", got)
    }
  })
}

// func TestBounds(t *testing.T){
  // t.Run("", func(t *testing.T){

// })
// }
