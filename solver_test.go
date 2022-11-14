package picross

import (
	"fmt"
	"testing"
)

func checkPicrWorker(t *testing.T, depth uint, clue []uint, input []CellState, expV []CellState) {
	w, e := NewPicrWorker(depth, clue, nil)
	if e != nil {
		t.Errorf(`PicrWorker creation failed: %v`, e)
		return
	}
	e = w.work(input)
	if e != nil {
		t.Errorf("unexpected error: %v", e)
		return
	}
	if !areSlicesEqual[CellState](w.getHint(), expV) {
		t.Errorf("unexpected value: expected %v, got %v", expV, w.getHint())
		return
	}
}

func TestPicrWorker(t *testing.T) {
	if _, e := NewPicrWorker(0, []uint{1}, nil); e == nil {
		t.Errorf(`unexpected success`)
	}
	checkPicrWorker(t, 3, []uint{}, []CellState{Any, Any, Any}, []CellState{Gap, Gap, Gap})
	checkPicrWorker(t, 3, []uint{3}, []CellState{Any, Any, Any}, []CellState{Fill, Fill, Fill})
	checkPicrWorker(t, 4, []uint{3}, []CellState{Any, Any, Any, Any}, []CellState{Any, Fill, Fill, Any})
	checkPicrWorker(t, 4, []uint{3}, []CellState{Gap, Any, Any, Any}, []CellState{Gap, Fill, Fill, Fill})
	checkPicrWorker(t, 4, []uint{3}, []CellState{Fill, Any, Any, Any}, []CellState{Fill, Fill, Fill, Gap})
	if w, _ := NewPicrWorker(4, []uint{3}, nil); w.work([]CellState{Any, Gap, Any, Any}) == nil {
		t.Errorf("unexpected success")
	}
}

func TestPicrAxis(t *testing.T) {
	a, e := NewPicrAxis(4, [][]uint{{3}}, nil)
	if e != nil {
		t.Fatalf(`%v`, e)
	}
	if e = a.work([][]CellState{{Gap, Any, Any, Any}}); e != nil {
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

func areSlices2Equal[K comparable](as [][]K, bs [][]K) bool {
	if len(as) != len(bs) {
		return false
	}
	for i, a := range as {
		if areSlicesEqual(a, bs[i]) {
			continue
		}
		return false
	}
	return true
}

func TestPicrTranspose(t *testing.T) {
	input := [][]CellState{[]CellState{Any, Fill, Gap}, []CellState{Fill, Gap, Any}}
	expected := [][]CellState{[]CellState{Any, Fill}, []CellState{Fill, Gap}, []CellState{Gap, Any}}
	got := picrTranspose(input)
	if !areSlices2Equal(expected, got) {
		t.Errorf("unexpected result: expected %v, got %v", expected, got)
	}
}

func TestPicrCountAny(t *testing.T) {
	input := [][]CellState{[]CellState{Any, Fill, Gap}, []CellState{Fill, Gap, Any}}
	expected := uint(2)
	got := picrCountAny(input)
	if expected != got {
		t.Errorf("unexpected result: expected %v, got %v", expected, got)
	}
}

func TestNewPicrSolver(t *testing.T) {
	if _, e := NewPicrSolver([][]uint{}, [][]uint{}, nil); e == nil {
		t.Fatalf("unexpected success")
	}
	if _, e := NewPicrSolver([][]uint{{1}}, [][]uint{}, nil); e == nil {
		t.Fatalf("unexpected success")
	}
	if _, e := NewPicrSolver([][]uint{}, [][]uint{{1}}, nil); e == nil {
		t.Fatalf("unexpected success")
	}
	if _, e := NewPicrSolver([][]uint{{1}}, [][]uint{{1}}, nil); e != nil {
		t.Fatalf("unexpected failure")
	}
}

func checkPicrSolverFail(t *testing.T, rowClues [][]uint, colClues [][]uint) {
	s, _ := NewPicrSolver(rowClues, colClues, nil)
	if s.solve() == nil {
		t.Fatalf(`unexpected success`)
	}
}

func TestPicrSolverSolveFail(t *testing.T) {
	checkPicrSolverFail(t, [][]uint{{0}}, [][]uint{{1}})
	checkPicrSolverFail(t, [][]uint{{1}}, [][]uint{{0}})
	checkPicrSolverFail(t, [][]uint{{1}}, [][]uint{{2}})
	checkPicrSolverFail(t, [][]uint{{2}}, [][]uint{{1}})
	checkPicrSolverFail(t, [][]uint{{1}, {1}}, [][]uint{{1}, {1}})
	checkPicrSolverFail(t, [][]uint{{1}, {2}}, [][]uint{{2}, {2}})
	checkPicrSolverFail(t, [][]uint{{2}, {1}}, [][]uint{{2}, {2}})
	checkPicrSolverFail(t, [][]uint{{2}, {2}}, [][]uint{{1}, {2}})
	checkPicrSolverFail(t, [][]uint{{2}, {2}}, [][]uint{{2}, {1}})
}

func checkPicrSolver(t *testing.T, rowClues [][]uint, colClues [][]uint, expOut [][]CellState) {
	solver, err := NewPicrSolver(rowClues, colClues, nil)
	if err != nil {
		t.Fatalf(`%v`, err)
	}
	err = solver.solve()
	if err != nil {
		t.Fatalf(`unexpected error for row clues %v and col clues %v: %v`, rowClues, colClues, err)
	}
	got := solver.getState()
	if !areSlices2Equal(expOut, got) {
		t.Fatalf(`result mismatch for row clues %v and col clues %v: expected %v, got %v`, rowClues, colClues, expOut, got)
	}
}

func printMap(m [][]CellState) {
	fmt.Println("")
	for _, row := range m {
		for _, v := range row {
			switch v {
			case Any:
				fmt.Printf("?")
			case Gap:
				fmt.Printf(".")
			case Fill:
				fmt.Printf("#")
			}
		}
		fmt.Println("")
	}
}

func str2Map(s string) [][]CellState {
	ans := make([][]CellState, 0)
	var row []CellState
	for _, c := range s {
		if row == nil {
			row = make([]CellState, 0)
		}
		if c == ' ' {
			continue
		}
		if c == '.' {
			row = append(row, Gap)
			continue
		}
		if c == '#' {
			row = append(row, Fill)
			continue
		}
		if c == '\n' {
			ans = append(ans, row)
			row = make([]CellState, 0)
			continue
		}
		panic(`unknown rune`)
	}
	ans = append(ans, row)
	return ans
}

func TestPicrSolver22Fill(t *testing.T) {
	// 2x2 Fill
	checkPicrSolver(t,
		[][]uint{{2}, {2}},
		[][]uint{{2}, {2}},
		str2Map(`##
                 ##`))
}

func TestPicrSolver55Horse(t *testing.T) {
	// 5x5 horse
	checkPicrSolver(t,
		[][]uint{{3}, {1, 1}, {4}, {3}, {1, 1}},
		[][]uint{{1}, {5}, {1, 2}, {3}, {2}},
		str2Map(`###..
                 .#..#
                 .####
                 .###.
                 .#.#.`))
}

func TestPicrSolverNotif(t *testing.T) {
	ch := make(chan PicrSolverNotification, 25)
	// 5x5 horse
	solver, _ := NewPicrSolver(
		[][]uint{{3}, {1, 1}, {4}, {3}, {1, 1}},
		[][]uint{{1}, {5}, {1, 2}, {3}, {2}},
		ch)
	solver.solve()
	got := make([]PicrSolverNotification, 0)
testPicrSolverNotifConsumeCh:
	for {
		select {
		case v := <-ch:
			got = append(got, v)
		default:
			break testPicrSolverNotifConsumeCh
		}
	}
	expected := []PicrSolverNotification{
		{1, 1, true}, {1, 2, true}, {1, 3, true}, {1, 4, false}, {1, 5, false},
		{2, 1, false}, {2, 2, true}, {2, 3, false}, {2, 4, false}, {2, 5, true},
		{3, 1, false}, {3, 2, true}, {3, 3, true}, {3, 4, true}, {3, 5, true},
		{4, 1, false}, {4, 2, true}, {4, 3, true}, {4, 4, true}, {4, 5, false},
		{5, 1, false}, {5, 2, true}, {5, 3, false}, {5, 4, true}, {5, 5, false},
	}
testPicrSolverNotifVerifyExpected:
	for _, v := range expected {
		for _, w := range got {
			if v == w {
				continue testPicrSolverNotifVerifyExpected
			}
		}
		t.Errorf(`expected notification not found: %v`, v)
	}
testPicrSolverNotifVerifyGot:
	for _, v := range got {
		for _, w := range expected {
			if v == w {
				continue testPicrSolverNotifVerifyGot
			}
		}
		t.Errorf(`unexpected notification found: %v`, v)
	}
}

func TestPicrSolverNonSquare(t *testing.T) {
	// A non-square puzzle
	checkPicrSolver(t,
		[][]uint{{3, 2, 2}, {3, 1, 1}, {4, 4}, {2, 1, 1, 1}},
		[][]uint{{4}, {4}, {3}, {2}, {1}, {1, 1}, {1, 1}, {1}, {1, 1}, {4}},
		str2Map(`###..##.##
                 ###.#....#
                 ####..####
                 ##.#.#...#`))
}

func TestPicrSolver1515Duck(t *testing.T) {
	// 15x15 duck
	checkPicrSolver(t,
		[][]uint{{3}, {5}, {4, 3}, {7}, {5}, {3}, {5}, {1, 8}, {3, 3, 3}, {7, 3, 2}, {5, 4, 2}, {8, 2}, {10}, {2, 3}, {6}},
		[][]uint{{3}, {4}, {5}, {4}, {5}, {6}, {3, 2, 1}, {2, 2, 5}, {4, 2, 6}, {8, 2, 3}, {8, 2, 1, 1}, {2, 6, 2, 1}, {4, 6}, {2, 4}, {1}},
		str2Map(`.........###...
                 ........#####..
                 .......####.###
                 .......#######.
                 ........#####..
                 .........###...
                 ........#####..
                 #.....########.
                 ###..###...###.
                 #######.###.##.
                 .#####.####.##.
                 .########..##..
                 ..##########...
                 ....##.###.....
                 ......######...`))
}

func TestPicrSolver2525Owl(t *testing.T) {
	// 25x25 owl
	checkPicrSolver(t,
		[][]uint{{3, 8, 3}, {23}, {3, 6, 3}, {7, 4, 7}, {2, 3, 2, 3, 4}, {2, 2, 2, 2, 2}, {1, 1, 1, 1, 2, 1, 1, 1, 2}, {1, 4, 2, 3, 4, 2}, {1, 3, 2, 1, 1, 2, 3, 1}, {1, 1, 2, 1, 1, 1, 1, 2, 1, 1}, {1, 2, 2, 1, 1, 1, 1, 2, 2, 1}, {1, 2, 2, 1, 1, 2, 2, 2}, {1, 4, 1, 1, 1, 4, 2}, {2, 1, 1, 2, 3, 1, 1, 1, 3}, {3, 2, 5, 1, 2}, {6, 5, 6}, {1, 3, 3, 1}, {3, 1, 2}, {3, 3}, {3, 3, 2, 4, 2}, {7, 11, 1}, {1, 5, 4, 1}, {1, 4, 4, 1}, {1, 17, 1}, {3, 5, 5, 3}},
		[][]uint{{1, 9, 2, 4}, {2, 3, 2, 2, 1}, {2, 2, 1, 1, 1, 3, 1, 1}, {3, 4, 2, 3, 1}, {3, 3, 3, 2, 2, 2}, {3, 1, 2, 1, 2, 1, 2}, {1, 2, 2, 2, 2, 2, 2, 2}, {1, 3, 2, 2, 2, 2, 2}, {2, 2, 4, 3, 3, 2}, {3, 2, 2, 2, 2}, {4, 5, 2, 1, 2}, {7, 4, 2, 1}, {8, 6, 3, 1}, {4, 5, 4, 2, 1}, {3, 3, 1, 2, 2, 2}, {2, 2, 4, 1, 3, 2}, {1, 2, 2, 2, 2, 3, 2}, {1, 2, 2, 2, 2, 1, 3, 2}, {3, 1, 2, 1, 1, 2, 2}, {3, 3, 3, 1, 1, 2}, {4, 4, 1, 1, 2}, {1, 2, 1, 1, 2, 2, 1}, {2, 3, 2, 2, 1, 1}, {2, 4, 3, 2, 1}, {1, 8, 3, 5}},
		str2Map(`###.....########......###
                 .#######################.
                 ...###...######...###....
                 .#######..####..#######..
                 .##...###..##..###..####.
                 ##.....##..##.##......##.
                 #...#.#..#.##.#..#.#...##
                 #...####.##.###.####...##
                 #.###..##.#..#.##..###..#
                 #..#.##.#.#..#.#.##.#...#
                 #.##.##.#.#..#.#.##.##..#
                 #..##..##.#..#.##..##..##
                 #...####.#..#.#.####...##
                 ##..#.#.##.###.#.#.#..###
                 .###...##.#####.#....##..
                 ...######.#####.######...
                 #...###....###..........#
                 ###.........#..........##
                 .###..................###
                 ..###.###...##.####..##..
                 ...#######.###########..#
                 #.......#####.####......#
                 #.####.............####.#
                 #...#################...#
                 ###...#####...#####...###`))
}

func TestPicrSolver4030Peacock(t *testing.T) {
	// 40x30 peacock
	// https://www.nonograms.org/nonograms/i/2098
	checkPicrSolver(t,
		[][]uint{{4, 7, 1}, {5, 9, 1, 2, 1}, {7, 12, 3, 2}, {9, 3, 7}, {12, 8, 1}, {6, 7, 2, 2, 2}, {4, 5, 6, 7}, {3, 7, 3, 4}, {6, 4, 2, 4}, {3, 3, 3, 4}, {3, 4, 6}, {3, 3, 1, 5, 6}, {2, 2, 4, 1, 2, 2, 6}, {5, 1, 1, 2, 3, 12}, {9, 1, 14}, {2, 2, 2, 9, 5}, {2, 2, 1, 1, 11, 4}, {1, 1, 12, 4}, {13, 3}, {3, 8, 4}, {4, 6, 4}, {5, 6}, {11, 9}, {25}, {26}, {24, 3, 2}, {21, 2, 2, 2}, {16, 1, 1, 1}, {4, 4, 3, 3}, {1, 2, 3}},
		[][]uint{{4, 1}, {4, 2}, {4, 1, 1, 2}, {5, 1, 2, 1, 1, 4}, {4, 2, 1, 4, 4}, {3, 2, 1, 3, 4}, {1, 5, 1, 3, 2, 4}, {1, 4, 2, 6, 3}, {2, 3, 2, 1, 3, 3}, {2, 3, 1, 3, 1, 2, 5}, {3, 5, 1, 4, 5}, {3, 9, 2, 6}, {3, 1, 4, 2, 6}, {2, 1, 2, 2, 5}, {2, 2, 1, 4, 6}, {2, 1, 2, 2, 2, 6}, {3, 1, 2, 3, 6}, {2, 1, 2, 3, 6}, {2, 2, 2, 1, 1, 6}, {1, 2, 1, 1, 7}, {1, 3, 7}, {1, 2, 7}, {3, 3, 4}, {1, 1, 3, 6}, {3, 5}, {4, 1, 3}, {7, 3}, {1, 7, 2, 1}, {2, 8, 4, 1}, {3, 8, 7}, {1, 2, 9, 5, 2}, {6, 8, 3, 1}, {8, 7, 8}, {3, 14, 6, 2}, {4, 7, 3, 5, 2}, {3, 7, 5, 1}, {1, 2, 13}, {2, 1, 1, 12}, {1, 1, 9}, {2, 1, 5}},
		str2Map(`####..#######........................#..
                 #####...#########.................#.##.#
                 #######...############...........###..##
                 #########.......###............#######..
                 ...############..............########..#
                 ......######..#######.......##.##....##.
                 ....####.#####....######...#######......
                 ...###....#######...###.......####......
                 .......######..####...##.......####.....
                 ......###..###...###............####....
                 .........###.####...............######..
                 ...###.###.#..#####..............######.
                 ..##..##.####.#.##............##.######.
                 ....#####.#.#.##.###........############
                 ..#########....#..........##############
                 ....##.##.##.............#########.#####
                 ...##.##.#.#............###########.####
                 ......#..#.............############.####
                 ......................#############.###.
                 .....................###..########.####.
                 ...................####..######...####..
                 .................#####..........######..
                 ..............###########...#########...
                 ...........#########################....
                 .........##########################.....
                 ...########################.###.##......
                 #####################..##....##.##......
                 .################......#.....#..#.......
                 ...####..####...............###.###.....
                 ...........................#..##.###....`))
}
