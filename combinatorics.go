package picross

// mapPermute returns a channel that provides all combinations of a single row (of a picross puzzle),
// `size` positions wide, that honor a `clue` of the run lenghts of the sequential marked pixels of the row.
// Each element of the answer is a bitmap of the row,
// where 'true' denotates a marked position
// and 'false' a gap.
func mapPermute(size uint, clue []uint) chan []bool {
	ans := make(chan []bool)
	go func() {
		defer close(ans)
		if len(clue) == 0 {
			ans <- make([]bool, size)
			return
		}
		picrs := picrPermute(size, clue)
		for picr := range picrs {
			ans <- picr2Map(size, picr)
		}
	}()
	return ans
}

// picr2Map converts gap-and-fill `lengths`
// (such as produced by picrPermute) to a bitmap representation.
// `size` is the amount of positions of the row.
func picr2Map(size uint, lengths []uint) []bool {
	var idx uint
	var pen bool
	ans := make([]bool, size)
	for _, length := range lengths {
		for i := uint(0); i < length; i++ {
			ans[idx] = pen
			idx += 1
		}
		pen = !pen
	}
	return ans
}

// picrPermute returns a channel that provides all combinations of a single row (of a picross puzzle),
// `size` positions wide, that honor a `clue` of the run lenghts of the sequential marked pixels of the row.
// Each element of the answer is a slice of run lengths representing:
// the first sequence of empty pixels,
// the first sequence of marked pixels,
// the next sequence of empty pixels,
// and so on.
func picrPermute(size uint, clue []uint) chan []uint {
    clueLen := uint(len(clue))
	var clueSum uint
	for _, v := range clue {
		clueSum += v
	}
	ans := make(chan []uint)
	if size < clueSum {
		close(ans)
		return ans
	}
	go func() {
        gapsLen := clueLen+1
        gapsSum := size-clueSum
        gaps := make([]uint, gapsLen)
        for i := range gaps {
            gaps[i] = 1
        }
        gaps[0] = 0
        gaps[gapsLen-1] = gapsSum+2-gapsLen
        for {
            e := make([]uint, clueLen+gapsLen)
            for i, v := range clue {
                e[1+2*i] = v
            }
            for i, v := range gaps {
                e[2*i] = v
            }
            ans <- e
            if !gapIterate(gaps) {
                break
            }
        }
		close(ans)
	}()
	return ans
}

// blend2 returns a slice whose first element is the first element of `as`,
// the second element is the first element of `bs`,
// the third element is the second element of `as` and so on.
func blend2(as []uint, bs []uint) []uint {
	lenb := len(bs)
	ans := make([]uint, len(as)+lenb)
	var idx int
	for i, a := range as {
		ans[idx] = a
		idx++
		if i < lenb {
			ans[idx] = bs[i]
			idx++
		}
	}
	for i := len(as); i < lenb; i++ {
		ans[idx] = bs[i]
		idx++
	}
	if idx != len(as)+len(bs) {
		panic("blend2: error: bad logic")
	}
	return ans
}

// gapCombine returns a channel that provides all combinations of integer numbers where
// the sum of all elements is `sum`,
// the amount of elements is `count`,
// the amount of elements is at least 2,
// the first and the last elements are equal or greater than 0,
// and the remaining elements are equal or greater than 1.
func gapCombine(sum uint, count uint) chan []uint {
	ans := make(chan []uint)
	if count < 2 || sum < count-2 {
		close(ans)
		return ans
	}
    buf := make([]uint, count)
    for i := range buf {
        buf[i] = 1
    }
    buf[0] = 0
    buf[count-1] = sum+2-count
    go func() {
        for {
            e := make([]uint, count)
            copy(e, buf)
            ans <- e
            if !gapIterate(buf) {
                close(ans)
                return
            }
        }
    }()
    return ans
}

func gapIterate(buf []uint) bool {
    bufLen := uint(len(buf))
    for i := int(bufLen)-1;; i-- {
        if i == 0 {
            return false
        }
        if i == int(bufLen)-1 {
            if buf[i] <= 0 {
                continue
            }
        } else {
            if buf[i] <= 1 {
                continue
            }
        }
        buf[i-1] += 1
        buf[i] -= 1
        for j := i; j < int(bufLen)-1; j++ {
            buf[j+1] += buf[j]-1
            buf[j] = 1
        }
        break
    }
    return true
}
