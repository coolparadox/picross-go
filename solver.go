package picross

import (
    "errors"
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

type PicrWorker struct {
    clue []uint
}

func NewPicrWorker(clue []uint) *PicrWorker {
    return &PicrWorker{clue: clue}
}

// emerge tries to detail a starting `hint` of the known state of a picross row (or column).
// The returned new state, when different than the input, contains less 'Any' values.
func (w *PicrWorker) emerge(hint []CellState) ([]CellState, error) {
    size := uint(len(hint))
    initialized := false
    pivot := make([]bool, size)
    dirty := make([]bool, size)
    permutations := mapPermute(size, w.clue)
emergeHintPermutations:
    for permutation := range permutations {
        for i, v := range permutation {
            oldCellState := hint[i]
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
        return []CellState{}, errors.New("unsolvable")
    }
    ans := make([]CellState, size)
    for i, v := range pivot {
        if dirty[i] {
            continue
        }
        if v {
            ans[i] = Fill
            continue
        }
        ans[i] = Gap
    }
    return ans, nil
}

