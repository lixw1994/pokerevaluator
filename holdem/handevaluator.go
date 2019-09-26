package holdem

// BitCount 计算handMask共表示几张牌
func BitCount(bitField uint64) int {
	return int(bits[(int)(bitField&0x00000000000000FF)] +
		bits[(int)((bitField&0x000000000000FF00)>>8)] +
		bits[(int)((bitField&0x0000000000FF0000)>>16)] +
		bits[(int)((bitField&0x00000000FF000000)>>24)] +
		bits[(int)((bitField&0x000000FF00000000)>>32)] +
		bits[(int)((bitField&0x0000FF0000000000)>>40)] +
		bits[(int)((bitField&0x00FF000000000000)>>48)] +
		bits[(int)((bitField&0xFF00000000000000)>>56)])
}

// EvaluateByMask 根据Card掩码计算牌力值
func EvaluateByMask(cardMask uint64, numberOfCards int) uint64 {
	retval := uint64(0)

	sc := ((cardMask >> (clubOffset)) & 0x1fff)
	sd := ((cardMask >> (diamondOffset)) & 0x1fff)
	sh := ((cardMask >> (heartOffset)) & 0x1fff)
	ss := ((cardMask >> (spadeOffset)) & 0x1fff)

	ranks := sc | sd | sh | ss
	nRanks := nBitsTable[ranks]               // 几张不重复的牌
	nDups := (uint64(numberOfCards) - nRanks) // 剩余几张牌

	if nRanks >= 5 {
		// 检测是否同花
		if nBitsTable[ss] >= 5 {
			if straightTable[ss] != 0 {
				return handtypeValueStraightFlush + (straightTable[ss] << topCardShift)
			}
			retval = handtypeValueFlush + (topFiveCardsTable[ss])
		} else if nBitsTable[sc] >= 5 {
			if straightTable[sc] != 0 {
				return handtypeValueStraightFlush + (straightTable[sc] << topCardShift)
			}
			retval = handtypeValueFlush + (topFiveCardsTable[sc])
		} else if nBitsTable[sd] >= 5 {
			if straightTable[sd] != 0 {
				return handtypeValueStraightFlush + (straightTable[sd] << topCardShift)
			}
			retval = handtypeValueFlush + (topFiveCardsTable[sd])
		} else if nBitsTable[sh] >= 5 {
			if straightTable[sh] != 0 {
				return handtypeValueStraightFlush + (straightTable[sh] << topCardShift)
			}
			retval = handtypeValueFlush + (topFiveCardsTable[sh])
		} else {
			st := straightTable[ranks]
			if st != 0 {
				retval = handtypeValueStraight + (st << topCardShift)
			}
		}
		// Another win -- if there can't be a FH/Quads (n_dups < 3),
		//           which is true most of the time when there is a made hand, then if we've
		//           found a five card hand, just return.  This skips the whole process of
		//           computing two_mask/three_mask/etc.
		if retval != 0 && nDups < 3 {
			return retval
		}
	}

	switch nDups {
	case 0:
		// It's a no-pair hand
		return handtypeValueHighCard + (topFiveCardsTable[ranks])
	case 1:
		// It's a one-pair hand
		twoMask := ranks ^ (sc ^ sd ^ sh ^ ss)

		retval = (handtypeValuePair + (topCardTable[twoMask] << topCardShift))
		t := ranks ^ twoMask
		kickers := (topFiveCardsTable[t] >> cardWidth) & (^fifthCardMask)
		retval += (kickers)
		return retval
	case 2:
		// Either two pair or trips
		twoMask := ranks ^ (sc ^ sd ^ sh ^ ss)
		if twoMask != 0 {
			t := ranks ^ twoMask
			retval = (handtypeValueTwoPair + (topFiveCardsTable[twoMask] & (topCardMask | secondCardMask)) + (topCardTable[t] << thirdCardShift))
			return retval
		}
		threeMask := ((sc & sd) | (sh & ss)) & ((sc & sh) | (sd & ss))
		retval = (handtypeValueTrips + (topCardTable[threeMask] << topCardShift))
		t := ranks ^ threeMask
		second := topCardTable[t]
		retval += (second << secondCardShift)
		t = t ^ (uint64(1) << second)
		retval += (topCardTable[t] << thirdCardShift)
	default:
		fourMask := sh & sd & sc & ss
		if fourMask != 0 {
			tc := topCardTable[fourMask]
			retval = (handtypeValueFourOfAKind + (tc << topCardShift) + ((topCardTable[ranks^(uint64(1)<<tc)]) << secondCardShift))
			return retval
		}

		twoMask := ranks ^ (sc ^ sd ^ sh ^ ss)
		if nBitsTable[twoMask] != nDups {
			threeMask := ((sc & sd) | (sh & ss)) & ((sc & sh) | (sd & ss))
			retval = handtypeValueFullHouse
			tc := topCardTable[threeMask]
			retval += (tc << topCardShift)
			t := (twoMask | threeMask) ^ (uint64(1) << tc)
			retval += (topCardTable[t] << secondCardShift)
			return retval
		}

		if retval != 0 { // flush and straight
			return retval
		}
		retval = handtypeValueTwoPair
		top := topCardTable[twoMask]
		retval += (top << topCardShift)
		second := topCardTable[twoMask^(uint64(1)<<top)]
		retval += (second << secondCardShift)
		retval += ((topCardTable[ranks^(uint64(1)<<top)^(uint64(1)<<second)]) << thirdCardShift)
		return retval
	}

	return retval
}

// Evaluate 计算牌力值 handValue
func Evaluate(cards []string) EvaluatorRes {
	cardMask := ParseHand(cards)
	return Evaluate7CardByMask(cardMask)
}

// EvaluatorRes EvaluateByMask计算结果
type EvaluatorRes struct {
	Value    uint64
	MaxHands []string
}

// Evaluate7CardByMask 根据Card掩码计算牌力值 只能是德扑7张牌
func Evaluate7CardByMask(cardMask uint64) (evRes EvaluatorRes) {
	numberOfCards := numberOfBoard + numberOfHand
	if BitCount(cardMask) != numberOfCards {
		return evRes
	}
	evRes.Value = uint64(0)

	sc := ((cardMask >> (clubOffset)) & 0x1fff)
	sd := ((cardMask >> (diamondOffset)) & 0x1fff)
	sh := ((cardMask >> (heartOffset)) & 0x1fff)
	ss := ((cardMask >> (spadeOffset)) & 0x1fff)

	ranks := sc | sd | sh | ss
	nRanks := nBitsTable[ranks]               // 几张不重复的牌
	nDups := (uint64(numberOfCards) - nRanks) // 剩余几张牌

	if nRanks >= 5 {
		// 检测是否同花
		if nBitsTable[ss] >= 5 {
			if straightTable[ss] != 0 {
				evRes.Value = handtypeValueStraightFlush + (straightTable[ss] << topCardShift)
				handIndexs := make([]uint64, 5)
				for i := 0; i < 5; i++ {
					rank := uint64(0)
					if straightTable[sc] >= uint64(i) {
						rank = straightTable[sc] - uint64(i)
					} else {
						rank = straightTable[sc] + 13 - uint64(i)
					}
					handIndexs[i] = rank + spadeOffset
				}
				evRes.MaxHands = cardsDesc(handIndexs)
				return evRes
			}
			evRes.Value = handtypeValueFlush + (topFiveCardsTable[ss])
		} else if nBitsTable[sc] >= 5 {
			if straightTable[sc] != 0 {
				evRes.Value = handtypeValueStraightFlush + (straightTable[sc] << topCardShift)
				handIndexs := make([]uint64, 5)
				for i := 0; i < 5; i++ {
					rank := uint64(0)
					if straightTable[sc] >= uint64(i) {
						rank = straightTable[sc] - uint64(i)
					} else {
						rank = straightTable[sc] + 13 - uint64(i)
					}
					handIndexs[i] = rank + clubOffset
				}
				evRes.MaxHands = cardsDesc(handIndexs)
				return evRes
			}
			evRes.Value = handtypeValueFlush + (topFiveCardsTable[sc])
		} else if nBitsTable[sd] >= 5 {
			if straightTable[sd] != 0 {
				evRes.Value = handtypeValueStraightFlush + (straightTable[sd] << topCardShift)
				handIndexs := make([]uint64, 5)
				for i := 0; i < 5; i++ {
					rank := uint64(0)
					if straightTable[sc] >= uint64(i) {
						rank = straightTable[sc] - uint64(i)
					} else {
						rank = straightTable[sc] + 13 - uint64(i)
					}
					handIndexs[i] = rank + diamondOffset
				}
				evRes.MaxHands = cardsDesc(handIndexs)
				return evRes
			}
			evRes.Value = handtypeValueFlush + (topFiveCardsTable[sd])
		} else if nBitsTable[sh] >= 5 {
			if straightTable[sh] != 0 {
				evRes.Value = handtypeValueStraightFlush + (straightTable[sh] << topCardShift)
				handIndexs := make([]uint64, 5)
				for i := 0; i < 5; i++ {
					rank := uint64(0)
					if straightTable[sc] >= uint64(i) {
						rank = straightTable[sc] - uint64(i)
					} else {
						rank = straightTable[sc] + 13 - uint64(i)
					}
					handIndexs[i] = rank + heartOffset
				}
				evRes.MaxHands = cardsDesc(handIndexs)
				return evRes
			}
			evRes.Value = handtypeValueFlush + (topFiveCardsTable[sh])
		} else {
			st := straightTable[ranks]
			if st != 0 {
				evRes.Value = handtypeValueStraight + (st << topCardShift)
				handIndexs := make([]uint64, 5)
				for i := 0; i < 5; i++ {
					rank := uint64(0)
					if st >= uint64(i) {
						rank = st - uint64(i)
					} else {
						rank = st + 13 - uint64(i)
					}
					handIndexs[i] = rank + whichOffset(rank, sc, sd, sh, ss)[0]
				}
				evRes.MaxHands = cardsDesc(handIndexs)
			}
		}
		// Another win -- if there can't be a FH/Quads (n_dups < 3),
		//           which is true most of the time when there is a made hand, then if we've
		//           found a five card hand, just return.  This skips the whole process of
		//           computing two_mask/three_mask/etc.
		if evRes.Value != 0 && nDups < 3 {
			return evRes
		}
	}

	switch nDups {
	case 0:
		// It's a no-pair hand
		fiveRankMask := topFiveCardsTable[ranks]
		evRes.Value = handtypeValueHighCard + fiveRankMask

		maxHandIndexs := make([]uint64, 5)
		for i := 4; fiveRankMask != 0; i-- {
			rank := fiveRankMask & oneCardMask
			fiveRankMask = (fiveRankMask >> cardWidth)
			maxHandIndexs[i] = rank + whichOffset(rank, sc, sd, sh, ss)[0]
		}
		evRes.MaxHands = cardsDesc(maxHandIndexs)
		return evRes
	case 1:
		// It's a one-pair hand
		twoMask := ranks ^ (sc ^ sd ^ sh ^ ss)

		evRes.Value = (handtypeValuePair + (topCardTable[twoMask] << topCardShift))
		t := ranks ^ twoMask
		kickers := (topFiveCardsTable[t] >> cardWidth) & (^fifthCardMask)
		evRes.Value += (kickers)

		maxHandIndexs := make([]uint64, 5)
		pairIndexs := whichOffset(topCardTable[twoMask], sc, sd, sh, ss)
		maxHandIndexs[0] = topCardTable[twoMask] + pairIndexs[0]
		maxHandIndexs[1] = topCardTable[twoMask] + pairIndexs[1]
		kickers = (kickers >> cardWidth)
		for i := 4; kickers != 0; i-- {
			rank := kickers & oneCardMask
			kickers = (kickers >> cardWidth)
			maxHandIndexs[i] = rank + whichOffset(rank, sc, sd, sh, ss)[0]
		}
		evRes.MaxHands = cardsDesc(maxHandIndexs)
		return evRes
	case 2:
		// Either two pair or trips
		twoMask := ranks ^ (sc ^ sd ^ sh ^ ss)
		if twoMask != 0 {
			t := ranks ^ twoMask
			evRes.Value = (handtypeValueTwoPair + (topFiveCardsTable[twoMask] & (topCardMask | secondCardMask)) + (topCardTable[t] << thirdCardShift))

			maxHandIndexs := make([]uint64, 5)
			twoPairCard := (topFiveCardsTable[twoMask] & (topCardMask | secondCardMask)) >> (cardWidth * 3)
			for i := 3; twoPairCard != 0; i = i - 2 {
				rank := twoPairCard & oneCardMask
				twoPairCard = (twoPairCard >> cardWidth)
				for index, offset := range whichOffset(rank, sc, sd, sh, ss) {
					maxHandIndexs[i-index] = rank + offset
				}
			}
			maxHandIndexs[4] = topCardTable[t] + whichOffset(topCardTable[t], sc, sd, sh, ss)[0]
			evRes.MaxHands = cardsDesc(maxHandIndexs)
			return evRes
		}
		threeMask := ((sc & sd) | (sh & ss)) & ((sc & sh) | (sd & ss))
		evRes.Value = (handtypeValueTrips + (topCardTable[threeMask] << topCardShift))
		leftFourMask := ranks ^ threeMask
		second := topCardTable[leftFourMask]
		evRes.Value += (second << secondCardShift)
		leftThreeMask := leftFourMask ^ (uint64(1) << second)
		third := topCardTable[leftThreeMask]
		evRes.Value += (third << thirdCardShift)

		maxHandIndexs := make([]uint64, 5)
		for index, offset := range whichOffset(topCardTable[threeMask], sc, sd, sh, ss) {
			maxHandIndexs[2-index] = topCardTable[threeMask] + offset
		}
		maxHandIndexs[3] = second + whichOffset(second, sc, sd, sh, ss)[0]
		maxHandIndexs[4] = third + whichOffset(third, sc, sd, sh, ss)[0]
		evRes.MaxHands = cardsDesc(maxHandIndexs)
	default:
		fourMask := sh & sd & sc & ss
		if fourMask != 0 {
			tc := topCardTable[fourMask]
			evRes.Value = (handtypeValueFourOfAKind + (tc << topCardShift) + ((topCardTable[ranks^(uint64(1)<<tc)]) << secondCardShift))

			maxHandIndexs := make([]uint64, 5)
			for index, offset := range whichOffset(tc, sc, sd, sh, ss) {
				maxHandIndexs[3-index] = tc + offset
			}
			maxHandIndexs[4] = topCardTable[ranks^(uint64(1)<<tc)] + whichOffset(topCardTable[ranks^(uint64(1)<<tc)], sc, sd, sh, ss)[0]
			evRes.MaxHands = cardsDesc(maxHandIndexs)
			return evRes
		}

		twoMask := ranks ^ (sc ^ sd ^ sh ^ ss)
		if nBitsTable[twoMask] != nDups {
			threeMask := ((sc & sd) | (sh & ss)) & ((sc & sh) | (sd & ss))
			evRes.Value = handtypeValueFullHouse
			tc := topCardTable[threeMask]
			evRes.Value += (tc << topCardShift)
			t := (twoMask | threeMask) ^ (uint64(1) << tc)
			evRes.Value += (topCardTable[t] << secondCardShift)

			maxHandIndexs := make([]uint64, 5)
			for index, offset := range whichOffset(tc, sc, sd, sh, ss) {
				maxHandIndexs[2-index] = tc + offset
			}
			for index, offset := range whichOffset(topCardTable[t], sc, sd, sh, ss) {
				maxHandIndexs[4-index] = topCardTable[t] + offset
			}
			evRes.MaxHands = cardsDesc(maxHandIndexs)
			return evRes
		}

		if evRes.Value != 0 { // flush and straight
			return evRes
		}

		evRes.Value = handtypeValueTwoPair
		top := topCardTable[twoMask]
		evRes.Value += (top << topCardShift)
		second := topCardTable[twoMask^(uint64(1)<<top)]
		evRes.Value += (second << secondCardShift)
		third := topCardTable[ranks^(uint64(1)<<top)^(uint64(1)<<second)]
		evRes.Value += (third << thirdCardShift)

		maxHandIndexs := make([]uint64, 5)
		for index, offset := range whichOffset(top, sc, sd, sh, ss) {
			maxHandIndexs[1-index] = top + offset
		}
		for index, offset := range whichOffset(second, sc, sd, sh, ss) {
			maxHandIndexs[3-index] = second + offset
		}
		maxHandIndexs[4] = third + whichOffset(third, sc, sd, sh, ss)[0]
		evRes.MaxHands = cardsDesc(maxHandIndexs)
		return evRes
	}

	return evRes
}

func whichOffset(cardRank, sc, sd, sh, ss uint64) (offsets []uint64) {
	target := uint64(1) << cardRank
	if target&sc != 0 {
		offsets = append(offsets, clubOffset)
		// fmt.Println(target, "clubOffset")
	}
	if target&sd != 0 {
		offsets = append(offsets, diamondOffset)
		// fmt.Println(target, "diamondOffset")
	}
	if target&sh != 0 {
		offsets = append(offsets, heartOffset)
		// fmt.Println(target, "heartOffset")
	}
	if target&ss != 0 {
		offsets = append(offsets, spadeOffset)
		// fmt.Println(target, "spadeOffset")
	}
	return offsets
}

func cardsDesc(handIndexs []uint64) []string {
	handDesc := make([]string, len(handIndexs))
	for i, handIndex := range handIndexs {
		handDesc[i] = cardTable[handIndex]
	}
	return handDesc
}
