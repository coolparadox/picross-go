package solver

func hfill(sum uint, count uint) chan []uint {
    ch := make(chan []uint)
    go func() {
        elm := make([]uint, 2)
        elm[0] = sum
        elm[1] = count
        ch <- elm
        close(ch)
    }()
    return ch
}

