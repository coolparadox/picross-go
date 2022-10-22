package picross

import (
    "testing"
)

func checkPicrWorker(t *testing.T, w *PicrWorker, input []CellState, expV []CellState) {
    e := w.work(input)
    if (e != nil) {
        t.Errorf("unexpected error: %v", e)
    }
    if !areSlicesEqual[CellState](w.getHint(), expV) {
        t.Errorf("unexpected value: expected %v, got %v", expV, w.getHint())
    }
}

func TestPicrWorker(t *testing.T) {
    checkPicrWorker(t, NewPicrWorker(4, []uint{3}), []CellState{Any, Any, Any, Any}, []CellState{Any, Fill, Fill, Any})
    checkPicrWorker(t, NewPicrWorker(4, []uint{3}), []CellState{Gap, Any, Any, Any}, []CellState{Gap, Fill, Fill, Fill})
    checkPicrWorker(t, NewPicrWorker(4, []uint{3}), []CellState{Fill, Any, Any, Any}, []CellState{Fill, Fill, Fill, Gap})
    if (NewPicrWorker(4, []uint{3}).work([]CellState{Any, Gap, Any, Any}) == nil) {
        t.Errorf("unexpected success")
    }
}

func TestPicrAxis(t *testing.T) {
    a := NewPicrAxis(4, [][]uint{{3}})
    if e := a.work([][]CellState{{Gap, Any, Any, Any}}); e != nil {
        t.Fatalf("unexpected error: %v", e)
    }
    got := a.getHint()
    if len(got) != 1 {
        t.Fatalf("length mismatch: expected 1, got %v", len(got))
    }
    expected := []CellState{Gap, Fill, Fill, Fill}
    if !areSlicesEqual[CellState](got[0], expected) {
        t.Errorf("value mismatch: expected %v, got %v", expected, got[0])
    }
}

