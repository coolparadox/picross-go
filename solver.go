package solver

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

