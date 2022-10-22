package picross

import (
    "errors"

    "golang.org/x/sync/errgroup"
)

type CellState uint

const (
    Any CellState = iota
    Gap
    Fill
)

func (c CellState) String() string {
	switch c {
	case Any:
		return "Any"
	case Gap:
		return "Gap"
	case Fill:
		return "Fill"
	}
    return "(error: unexpected)"
}

// PicrWorker handles a single row (column) of a picross puzzle.
type PicrWorker struct {
    clue []uint
    hint []CellState
}

func NewPicrWorker(depth uint, clue []uint) *PicrWorker {
    return &PicrWorker{clue:clue, hint:make([]CellState, depth)}
}

func (w *PicrWorker) getHint() []CellState {
    return w.hint
}

// work tries to detail a starting `hint` of the known state of a picross row (or column).
// The returned new state, when different than the input, contains less 'Any' values.
func (w *PicrWorker) work(hint []CellState) error {
    if len(hint) != len(w.hint) {
        panic("mismatched hint length")
    }
    for i, v := range hint {
        if v == Any || w.hint[i] == Any {
            continue
        }
        if v != w.hint[i] {
            return errors.New("nonsense hint")
        }
    }
    for i, v := range hint {
        if v == Any {
            continue
        }
        w.hint[i] = v
    }
    size := uint(len(w.hint))
    initialized := false
    pivot := make([]bool, size)
    dirty := make([]bool, size)
    permutations := mapPermute(size, w.clue)
emergeHintPermutations:
    for permutation := range permutations {
        for i, v := range permutation {
            oldCellState := w.hint[i]
            if (v && oldCellState == Gap) || (!v && oldCellState == Fill) {
                continue emergeHintPermutations
            }
        }
        if !initialized {
            initialized = true
            copy(pivot, permutation)
            continue emergeHintPermutations
        }
        for i, v := range permutation {
            if v != pivot[i] {
                dirty[i] = true
            }
        }
    }
    if !initialized {
        return errors.New("no solution")
    }
    for i, v := range pivot {
        if dirty[i] {
            continue
        }
        if v {
            w.hint[i] = Fill
            continue
        }
        w.hint[i] = Gap
    }
    return nil
}

type PicrAxis struct {
    workers []*PicrWorker
}

func NewPicrAxis(depth uint, clues[][]uint) *PicrAxis {
    workers := make([]*PicrWorker, len(clues))
    for i, clue := range clues {
        workers[i] = NewPicrWorker(depth, clue)
    }
    return &PicrAxis{workers:workers}
}

func (a *PicrAxis) getHint() [][]CellState {
    ans := make([][]CellState, len(a.workers))
    for i, w := range a.workers {
        ans[i] = w.getHint()
    }
    return ans
}

func (a *PicrAxis) work(hint [][]CellState) error {
    if len(hint) != len(a.workers) {
        panic("hint length mismatch")
    }
    g := new(errgroup.Group)
    for i, w := range a.workers {
        i := i
        w := w
        g.Go(func() error { return w.work(hint[i]) })
    }
    return g.Wait()
}

