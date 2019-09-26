package holdem

func comb(n, m int, emit func([]int)) {
	s := make([]int, m)
	last := m - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			s[i] = j
			if i == last {
				emit(s)
			} else {
				rc(i+1, j+1)
			}
		}
		return
	}
	rc(0, 0)
}

// CardComb 计算所有牌组合
func CardComb(n int, deadCardsMask uint64) []uint64 {
	res := make([]uint64, 0, numberOfGroups)
	if n <= 0 {
		return res
	}

	comb(numberOfPoker, n, func(one []int) {
		combMask := uint64(0)
		for i := 0; i < n; i++ {
			combMask = combMask | cardMasksTable[one[i]]
		}
		if (combMask & deadCardsMask) == 0 {
			res = append(res, combMask)
		}
	})

	return res
}
