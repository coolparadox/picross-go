package picross

import (
    "testing"
)

func checkPicrWorker(t *testing.T, w *PicrWorker, input []CellState, expV []CellState) {
    v, e := w.emerge(input)
    if (e != nil) {
        t.Errorf("unexpected error: %v", e)
    }
    if !areSlicesEqual[CellState](v, expV) {
        t.Errorf("unexpected value: expected %v, got %v", expV, v)
    }
}

func TestEmergeHint(t *testing.T) {
    w := NewPicrWorker([]uint{3})
    checkPicrWorker(t, w, []CellState{Any, Any, Any, Any}, []CellState{Any, Fill, Fill, Any})
    checkPicrWorker(t, w, []CellState{Gap, Any, Any, Any}, []CellState{Gap, Fill, Fill, Fill})
    checkPicrWorker(t, w, []CellState{Fill, Any, Any, Any}, []CellState{Fill, Fill, Fill, Gap})
    _, e := w.emerge([]CellState{Any, Gap, Any, Any})
    if (e == nil) {
        t.Errorf("unexpected success")
    }
}

