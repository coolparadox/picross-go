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

func TestHFill0(t *testing.T) {
    exps := [][]uint{
        {1,1,3},
        {1,2,2},
        {1,3,1},
        {2,1,2},
        {2,2,1},
        {3,1,1},
    }
    fills_ch := HFill(5, 3)
    fills := [][]uint{}
    for fill := range fills_ch {
        fills = append(fills, fill)
    }
    if len(fills) != len(exps) {
        t.Fatalf(`length mismatch: expected %v, got %v`, len(exps), len(fills))
    }
    for _, exp := range exps {
        if !uint_slice_in_slice(exp, fills) {
            t.Errorf(`expected element not found: %v`, exp)
        }
    }
    for _, fill := range fills {
        if !uint_slice_in_slice(fill, exps) {
            t.Errorf(`unexpected element found: %v`, fill)
        }
    }
}

