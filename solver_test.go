package solver

import (
    "testing"
)

func TestHFill0(t *testing.T) {
    ch := hfill(3, 5)
    var count uint = 0
    for elm := range ch {
        count++
        if count > 1 {
            t.Fatalf(`len fail`)
        }
        if elm[0] != 3 {
            t.Fatalf(`0 fail`)
        }
        if elm[1] != 5 {
            t.Fatalf(`0 fail`)
        }
    }
    if count < 1 {
        t.Fatalf(`len fail`)
    }
}
