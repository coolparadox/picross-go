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

// gapIterate takes a slice of run lenghts and in-place rearranges it so that:
// - it has the same amount of elements,
// - it has the same sum of elements,
// - its first and last elements are equal or greater than zero,
// - its remaining elements are equal or greater than one,
// - it's the next realization from a sorted set of all combinations that honor the previous invariants.
// Returns true on success achieving a new combination, of false otherwise.
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
