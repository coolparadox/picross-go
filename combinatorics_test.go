package picross

import (
	"testing"
)

func areSlicesEqual[K comparable](as []K, bs []K) bool {
	if len(as) != len(bs) {
		return false
	}
	for i, a := range as {
		if a != bs[i] {
			return false
		}
	}
	return true
}

func isSliceInSlice[K comparable](elm []K, source [][]K) bool {
	for _, e := range source {
		if areSlicesEqual[K](e, elm) {
			return true
		}
	}
	return false
}

func checkExpectedSlices[K comparable](t *testing.T, ch chan []K, exps [][]K) {
	fills := [][]K{}
	for fill := range ch {
		if !isSliceInSlice[K](fill, exps) {
			t.Errorf(`unexpected element found: %v`, fill)
		}
		fills = append(fills, fill)
	}
	for _, exp := range exps {
		if !isSliceInSlice[K](exp, fills) {
			t.Errorf(`expected element not found: %v`, exp)
		}
	}
}

func TestGapCombine51(t *testing.T) {
	exps := [][]uint{}
	checkExpectedSlices[uint](t, gapCombine(5, 1), exps)
}

func TestGapCombine52(t *testing.T) {
	exps := [][]uint{
		{0, 5},
		{1, 4},
		{2, 3},
		{3, 2},
		{4, 1},
		{5, 0},
	}
	checkExpectedSlices[uint](t, gapCombine(5, 2), exps)
}

func TestGapCombine53(t *testing.T) {
	exps := [][]uint{
		{0, 5, 0},
		{0, 4, 1},
		{0, 3, 2},
		{0, 2, 3},
		{0, 1, 4},
		{1, 4, 0},
		{1, 3, 1},
		{1, 2, 2},
		{1, 1, 3},
		{2, 3, 0},
		{2, 2, 1},
		{2, 1, 2},
		{3, 2, 0},
		{3, 1, 1},
		{4, 1, 0},
	}
	checkExpectedSlices[uint](t, gapCombine(5, 3), exps)
}

func TestBlend2(t *testing.T) {
	if !areSlicesEqual[uint](blend2([]uint{10, 11}, []uint{20, 21}), []uint{10, 20, 11, 21}) {
		t.Errorf(`mismatch`)
	}
	if !areSlicesEqual[uint](blend2([]uint{10}, []uint{20, 21}), []uint{10, 20, 21}) {
		t.Errorf(`mismatch`)
	}
	if !areSlicesEqual[uint](blend2([]uint{10, 11}, []uint{20}), []uint{10, 20, 11}) {
		t.Errorf(`mismatch`)
	}
}

func TestPicrPermute(t *testing.T) {
	exps := [][]uint{
		{0, 2, 5, 3, 0},
		{0, 2, 4, 3, 1},
		{0, 2, 3, 3, 2},
		{0, 2, 2, 3, 3},
		{0, 2, 1, 3, 4},
		{1, 2, 4, 3, 0},
		{1, 2, 3, 3, 1},
		{1, 2, 2, 3, 2},
		{1, 2, 1, 3, 3},
		{2, 2, 3, 3, 0},
		{2, 2, 2, 3, 1},
		{2, 2, 1, 3, 2},
		{3, 2, 2, 3, 0},
		{3, 2, 1, 3, 1},
		{4, 2, 1, 3, 0},
	}
	checkExpectedSlices[uint](t, picrPermute(10, []uint{2, 3}), exps)
}

func TestPicr2Map(t *testing.T) {
	expected := []bool{false, true, true, true, false, false, true, true}
	got := picr2Map(8, []uint{1, 3, 2, 2, 0})
	if !areSlicesEqual[bool](got, expected) {
		t.Errorf(`mismatch: got %v, expected %v`, got, expected)
	}
}

func TestMapPermute(t *testing.T) {
	exps := [][]bool{
		{true, true, false, false, true},
		{true, true, false, true, false},
		{false, true, true, false, true},
	}
	checkExpectedSlices[bool](t, mapPermute(5, []uint{2, 1}), exps)
}
