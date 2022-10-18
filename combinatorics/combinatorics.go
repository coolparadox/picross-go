// Package combinatorics offers free functions for brute-force approaches towards solving picross
// puzzles.
package combinatorics

// xFill returns a channel that provides all combinations of integer numbers where
// the sum of all elements is `sum`,
// the amount of elements is `count`,
// the amount of elements is at least 2,
// the first and the last elements are equal or greater than 0,
// and the remaining elements are equal or greater than 1.
func xFill(sum uint, count uint) chan []uint {
    ans := make(chan []uint)
    if count < 2 || sum < count-2 {
        close(ans)
        return ans
    }
    if count == 2 {
        go func() {
            for head := uint(0); head <= sum; head++ {
                elm := make([]uint, 2)
                elm[0] = head
                elm[1] = sum-head
                ans <- elm
            }
            close(ans)
        }()
        return ans
    }
    go func() {
        for head := uint(0); head <= sum-(count-2); head++ {
            for last := uint(0); last <= sum-(count-2)-head; last++ {
                middles := hFill(sum-head-last, count-2)
                for middle := range middles {
                    ans <- append(append([]uint{head}, middle...), last)
                }
            }
        }
        close(ans)
    }()
    return ans
}

// hFill returns a channel that provides all combinations of integer numbers where
// the sum of all elements is `sum`,
// the amount of elements is `count`,
// and each element is equal or greater than 1.
func hFill(sum uint, count uint) chan []uint {
    ans := make(chan []uint)
    if count < 1 || sum < count {
        close(ans)
        return ans
    }
    go func() {
        if count == 1 {
            elm := make([]uint, 1)
            elm[0] = sum
            ans <- elm
            close(ans)
            return
        }
        for last := uint(1); last <= sum-(count-1) ; last++ {
            predecessors := hFill(sum-last, count-1)
            for predecessor := range predecessors {
                ans <- append(predecessor, last)
            }
        }
        close(ans)
    }()
    return ans
}

