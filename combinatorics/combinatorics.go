// Package combinatorics offers free functions for brute-force approaches towards solving picross
// puzzles.
package combinatorics

// HFill returns a channel that provides all combinations of integer numbers where
// the sum of all elements is `sum`,
// the amount of elements is `count`,
// and each element is equal or greater than 1.
func HFill(sum uint, count uint) chan []uint {
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
            predecessors := HFill(sum-last, count-1)
            for predecessor := range predecessors {
                ans <- append(predecessor, last)
            }
        }
        close(ans)
    }()
    return ans
}

