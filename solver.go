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

type PicrWorkerNotification struct {
	position uint
	value    CellState
}

// PicrWorker handles a single row (column) of a picross puzzle.
type PicrWorker struct {
	isPrimed bool
	clue     []uint
	hint     []CellState
	notifCh  chan PicrWorkerNotification
}

func NewPicrWorker(depth uint, clue []uint, notifCh chan PicrWorkerNotification) (*PicrWorker, error) {
	if depth < 1 {
		return nil, errors.New("PicrWorker: zero depth")
	}
	return &PicrWorker{clue: clue, hint: make([]CellState, depth), notifCh: notifCh}, nil
}

func (w *PicrWorker) getHint() []CellState {
	return w.hint
}

func (w *PicrWorker) getNotifCh() chan PicrWorkerNotification {
	return w.notifCh
}

// work tries to detail a starting `hint` of the known state of a picross row (or column).
// The updated new state, when different than the input, contains less 'Any' values.
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
			if w.notifCh != nil {
				w.notifCh <- PicrWorkerNotification{position: uint(i), value: v}
			}
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
		oldHint := w.hint[i]
		if v {
			w.hint[i] = Fill
		} else {
			w.hint[i] = Gap
		}
		if w.notifCh != nil && w.hint[i] != oldHint {
			w.notifCh <- PicrWorkerNotification{position: uint(i), value: w.hint[i]}
		}
	}
	return nil
}

type PicrAxisNotification struct {
	workerIdx uint
	workerPos uint
	value     CellState
}

type PicrAxis struct {
	workers []*PicrWorker
	notifCh chan PicrAxisNotification
}

func NewPicrAxis(depth uint, clues [][]uint, notifCh chan PicrAxisNotification) (*PicrAxis, error) {
	if len(clues) < 1 {
		return nil, errors.New("PicrAxis: empty clues")
	}
	workers := make([]*PicrWorker, len(clues))
	for i, clue := range clues {
		var err error
		var workerNotif chan PicrWorkerNotification
		if notifCh != nil {
			workerNotif = make(chan PicrWorkerNotification, depth)
		}
		workers[i], err = NewPicrWorker(depth, clue, workerNotif)
		if err != nil {
			return nil, err
		}
	}
	return &PicrAxis{workers: workers, notifCh: notifCh}, nil
}

func (a *PicrAxis) getHint() [][]CellState {
	ans := make([][]CellState, len(a.workers))
	for i, w := range a.workers {
		ans[i] = w.getHint()
	}
	return ans
}

func (a *PicrAxis) getNotifCh() chan PicrAxisNotification {
	return a.notifCh
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
	if err := g.Wait(); err != nil {
		return err
	}
	if a.notifCh != nil {
		for i, w := range a.workers {
			ch := w.getNotifCh()
		picrAxisWorkConsumeWorkerChan:
			for {
				select {
				case n := <-ch:
					a.notifCh <- PicrAxisNotification{workerIdx: uint(i), workerPos: n.position, value: n.value}
				default:
					break picrAxisWorkConsumeWorkerChan
				}
			}
		}
	}
	return nil
}

type PicrSolverNotification struct {
	row  uint
	col  uint
	mark bool
}

type PicrSolver struct {
	row     *PicrAxis
	col     *PicrAxis
	notifCh chan PicrSolverNotification
}

func NewPicrSolver(rowClues [][]uint, colClues [][]uint, notifCh chan PicrSolverNotification) (*PicrSolver, error) {
	var axisNotif chan PicrAxisNotification
	if notifCh != nil {
		axisNotif = make(chan PicrAxisNotification, len(rowClues)*len(colClues))
	}
	row, err := NewPicrAxis(uint(len(colClues)), rowClues, axisNotif)
	if err != nil {
		return nil, err
	}
	col, err := NewPicrAxis(uint(len(rowClues)), colClues, nil)
	if err != nil {
		return nil, err
	}
	return &PicrSolver{row: row, col: col, notifCh: notifCh}, nil
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
		if s.notifCh != nil {
			axisNotifCh := s.row.getNotifCh()
		picrSolverSolveConsumeAxisNotifCh:
			for {
				select {
				case axisNotif := <-axisNotifCh:
					s.notifCh <- PicrSolverNotification{row: axisNotif.workerIdx + 1, col: axisNotif.workerPos + 1, mark: axisNotif.value == Fill}
				default:
					break picrSolverSolveConsumeAxisNotifCh
				}
			}
		}
		n := picrCountAny(s.row.getHint())
		if n == n_unknown {
			return errors.New("PicrSolver: dubious puzzle")
		}
		n_unknown = n
	}
	return s.col.work(picrTranspose(s.row.getHint()))
}
