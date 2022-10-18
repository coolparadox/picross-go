package combinatorics

import (
    "testing"
)

func uint_slices_equal(as []uint, bs []uint) bool {
    if len(as) != len(bs) {
        return false
    }
    for i, a := range(as) {
        if a != bs[i] {
            return false
        }
    }
    return true
}

func uint_slice_in_slice(elm []uint, source [][]uint) bool {
    for _, e := range source {
        if uint_slices_equal(e, elm) {
            return true
        }
    }
    return false
}

func checkFills(t *testing.T, ch chan []uint, exps [][]uint) {
    fills := [][]uint{}
    for fill := range ch {
        if !uint_slice_in_slice(fill, exps) {
            t.Errorf(`unexpected element found: %v`, fill)
        }
        fills = append(fills, fill)
    }
    for _, exp := range exps {
        if !uint_slice_in_slice(exp, fills) {
            t.Errorf(`expected element not found: %v`, exp)
        }
    }
}

func TestHFill51(t *testing.T) {
    exps := [][]uint{
        {5},
    }
    checkFills(t, hFill(5, 1), exps)
}

func TestHFill53(t *testing.T) {
    exps := [][]uint{
        {1, 1, 3},
        {1, 2, 2},
        {1, 3, 1},
        {2, 1, 2},
        {2, 2, 1},
        {3, 1, 1},
    }
    checkFills(t, hFill(5, 3), exps)
}

func TestXFill51(t *testing.T) {
    exps := [][]uint{}
    checkFills(t, xFill(5, 1), exps)
}

func TestXFill52(t *testing.T) {
    exps := [][]uint{
        {0, 5},
        {1, 4},
        {2, 3},
        {3, 2},
        {4, 1},
        {5, 0},
    }
    checkFills(t, xFill(5, 2), exps)
}

func TestXFill53(t *testing.T) {
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
    checkFills(t, xFill(5, 3), exps)
}

func TestBlend2(t *testing.T) {
    if !uint_slices_equal(blend2([]uint{10,11}, []uint{20,21}), []uint{10,20,11,21}) {
        t.Errorf(`mismatch`)
    }
    if !uint_slices_equal(blend2([]uint{10}, []uint{20,21}), []uint{10,20,21}) {
        t.Errorf(`mismatch`)
    }
    if !uint_slices_equal(blend2([]uint{10,11}, []uint{20}), []uint{10,20,11}) {
        t.Errorf(`mismatch`)
    }
}

