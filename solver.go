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
	isPrimed bool
	clue     []uint
	hint     []CellState
}

func NewPicrWorker(depth uint, clue []uint) (*PicrWorker, error) {
	if depth < 1 {
		return nil, errors.New("PicrWorker: zero depth")
	}
	return &PicrWorker{clue: clue, hint: make([]CellState, depth)}, nil
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
			return errors.New("PicrWorker: nonsense hint")
		}
	}
	anyChange := false
	for i, v := range hint {
		if v == Any {
			continue
		}
		if w.hint[i] != v {
			anyChange = true
		}
		w.hint[i] = v
	}
	if !anyChange && w.isPrimed {
		return nil
	}
	w.isPrimed = true
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
		return errors.New("PicrWorker: no solution")
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

func NewPicrAxis(depth uint, clues [][]uint) (*PicrAxis, error) {
	if len(clues) < 1 {
		return nil, errors.New("PicrAxis: empty clues")
	}
	workers := make([]*PicrWorker, len(clues))
	for i, clue := range clues {
		var err error
		workers[i], err = NewPicrWorker(depth, clue)
		if err != nil {
			return nil, err
		}
	}
	return &PicrAxis{workers: workers}, nil
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

type PicrSolver struct {
	row *PicrAxis
	col *PicrAxis
}

func NewPicrSolver(rowClues [][]uint, colClues [][]uint) (*PicrSolver, error) {
	row, err := NewPicrAxis(uint(len(colClues)), rowClues)
	if err != nil {
		return nil, err
	}
	col, err := NewPicrAxis(uint(len(rowClues)), colClues)
	if err != nil {
		return nil, err
	}
	return &PicrSolver{row: row, col: col}, nil
}

func (s *PicrSolver) getState() [][]CellState {
	return s.row.getHint()
}

func picrTranspose(mat [][]CellState) [][]CellState {
	ans := make([][]CellState, 0)
	for rowIdx := range mat[0] {
		col := make([]CellState, 0)
		for _, row := range mat {
			col = append(col, row[rowIdx])
		}
		ans = append(ans, col)
	}
	return ans
}

func picrCountAny(mat [][]CellState) uint {
	var ans uint
	for _, row := range mat {
		for _, elm := range row {
			if elm != Any {
				continue
			}
			ans += 1
		}
	}
	return ans
}

func (s *PicrSolver) solve() error {
	n_unknown := picrCountAny(s.row.getHint())
	for n_unknown > 0 {
		if err := s.col.work(picrTranspose(s.row.getHint())); err != nil {
			return err
		}
		if err := s.row.work(picrTranspose(s.col.getHint())); err != nil {
			return err
		}
		n := picrCountAny(s.row.getHint())
		if n == n_unknown {
			return errors.New("PicrSolver: dubious puzzle")
		}
		n_unknown = n
	}
	return s.col.work(picrTranspose(s.row.getHint())) // FIXME: optimize: just check validity instead of working
}
