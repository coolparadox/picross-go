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
		gaps_ch := xFill(size-clueSum, uint(len(clue))+1)
		for gaps := range gaps_ch {
			ans <- append(gaps[:1], blend2(clue, gaps[1:])...)
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
				elm[1] = sum - head
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
// the amount of elements is `size`,
// and each element is equal or greater than 1.
func hFill(sum uint, size uint) chan []uint {
	ans := make(chan []uint)
	if size < 1 || sum < size {
		close(ans)
		return ans
	}
    buf := make([]uint, size)
    var it func(uint, uint)
    it = func(sum uint, idx uint) {
        if idx == size-1 {
            buf[idx] = sum
            elem := make([]uint, size)
            copy(elem, buf)
            ans <- elem
            return
        }
        for v := uint(1); v <= sum+idx+1-size; v++ {
            v := v
            buf[idx] = v
            it(sum-v, idx+1)
        }
    }
    go func() {
        it(sum, 0)
        close(ans)
    }()
    return ans
}
