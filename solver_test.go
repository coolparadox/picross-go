package solver

import (
    "testing"
)

func TestHFill0(t *testing.T) {
    got := hfill(3, 5)
    if len(got) != 2 {
        t.Fatalf(`len fail`)
    }
    if got[0] != 3 {
        t.Fatalf(`0 fail`)
    }
    if got[1] != 5 {
        t.Fatalf(`0 fail`)
    }
}
