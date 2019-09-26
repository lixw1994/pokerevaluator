package holdem

import (
	"fmt"
	"time"
)

// WinStat 玩家胜利统计
type WinStat struct {
	WinCounts []uint64
	TieCounts []uint64
	Count     uint64
}

// CalcWithoutOpponentLimit 限制次数的模拟 - 一般用在非all-allin场景
func CalcWithoutOpponentLimit(myhandMask uint64, boardMask uint64, opponentCount int, limit int) WinStat {
	totalWinStat := WinStat{make([]uint64, opponentCount+1), make([]uint64, opponentCount+1), uint64(0)}
	deadCardMask := (myhandMask | boardMask)
	allHands := make([]uint64, opponentCount+1)
	allHands[0] = myhandMask
	boardCount := BitCount(boardMask)
	for i := 0; i < limit; i++ {
		masks := RandomNCardMask(opponentCount*numberOfHand+numberOfBoard-boardCount, deadCardMask)

		for oppI := 0; oppI < opponentCount; oppI++ {
			allHands[oppI+1] = (masks[2*oppI] | masks[2*oppI+1])
		}
		extBoardMask := uint64(0)
		for bI := 0; bI < numberOfBoard-boardCount; bI++ {
			extBoardMask |= masks[2*opponentCount+bI]
		}
		totalWinStat.mergeWinStat(CalcPlayersByMask(allHands, boardMask|extBoardMask))
	}

	return totalWinStat
}

// CalcWithOpponentLimit 限制次数的模拟 - 用在all-allin场景
func CalcWithOpponentLimit(handMasks []uint64, boardMask uint64, limit int) WinStat {
	boardCount := BitCount(boardMask)
	if boardCount >= 5 {
		return CalcPlayersByMask(handMasks, boardMask)
	}

	opponentCount := len(handMasks) - 1
	totalWinStat := WinStat{make([]uint64, opponentCount+1), make([]uint64, opponentCount+1), uint64(0)}
	deadCardsMask := boardMask
	for _, handMask := range handMasks {
		deadCardsMask = deadCardsMask | handMask
	}
	boardCombs := CardComb(numberOfBoard-boardCount, deadCardsMask)
	for index, boardComb := range boardCombs {
		totalWinStat.mergeWinStat(CalcPlayersByMask(handMasks, boardComb|boardMask))
		if index >= limit {
			break
		}
	}
	return totalWinStat
}

// CalcWithOpponent 不限制次数的模拟 - 用在all-allin场景
func CalcWithOpponent(handMasks []uint64, boardMask uint64) WinStat {
	boardCount := BitCount(boardMask)
	if boardCount >= 5 {
		return CalcPlayersByMask(handMasks, boardMask)
	}

	opponentCount := len(handMasks) - 1
	totalWinStat := WinStat{make([]uint64, opponentCount+1), make([]uint64, opponentCount+1), uint64(0)}
	deadCardsMask := boardMask
	for _, handMask := range handMasks {
		deadCardsMask = deadCardsMask | handMask
	}
	boardCombs := CardComb(numberOfBoard-boardCount, deadCardsMask)
	for _, boardComb := range boardCombs {
		totalWinStat.mergeWinStat(CalcPlayersByMask(handMasks, boardComb|boardMask))
	}
	return totalWinStat
}

// CalcWithOneOpponent 一个对手所有情况 ~2min
func CalcWithOneOpponent(myhandMask uint64, boardMask uint64) WinStat {
	t := time.Now()

	totalWinStat := WinStat{make([]uint64, 2), make([]uint64, 2), uint64(0)}

	hands := CardComb(numberOfHand, myhandMask|boardMask)
	for _, opponentMask := range hands {
		totalWinStat.mergeWinStat(CalcPlayersByMask([]uint64{myhandMask, opponentMask}, boardMask))
	}

	for i, winCount := range totalWinStat.WinCounts {
		fmt.Printf("player-%v win rate %v tie rate %v \n", i, float64(winCount)/float64(totalWinStat.Count), float64(totalWinStat.TieCounts[i])/float64(totalWinStat.Count))
		fmt.Printf("player-%v win rate %v tie rate %v \n", i, float64(winCount)/float64(totalWinStat.Count), float64(totalWinStat.TieCounts[i])/float64(totalWinStat.Count))
	}

	elapsed := time.Since(t)
	fmt.Println("CalcOnePlayer elapsed: ", elapsed)

	return totalWinStat
}

// CalcWithTwoOpponent 两个未知对手的情况 ~2000min
func CalcWithTwoOpponent(myhandMask uint64, boardMask uint64) WinStat {
	t := time.Now()

	totalWinStat := WinStat{make([]uint64, 3), make([]uint64, 3), uint64(0)}

	handOnes := CardComb(numberOfHand, myhandMask|boardMask)
	for _, opponentOneMask := range handOnes {
		handTwos := CardComb(numberOfHand, myhandMask|opponentOneMask|boardMask)
		for _, opponentTwoMask := range handTwos {
			totalWinStat.mergeWinStat(CalcPlayersByMask([]uint64{myhandMask, opponentOneMask, opponentTwoMask}, boardMask))
		}
	}

	for i, winCount := range totalWinStat.WinCounts {
		fmt.Printf("player-%v win rate %v tie rate %v \n", i, float64(winCount)/float64(totalWinStat.Count), float64(totalWinStat.TieCounts[i])/float64(totalWinStat.Count))
	}

	elapsed := time.Since(t)
	fmt.Println("CalcOnePlayer elapsed: ", elapsed)

	return totalWinStat
}

// CalcPlayersByMask 计算玩家胜负情况
func CalcPlayersByMask(handMasks []uint64, boardMask uint64) WinStat {
	winCounts := make([]uint64, len(handMasks))
	tieCounts := make([]uint64, len(handMasks))
	count := uint64(0)

	handValues := make([]uint64, 0, len(handMasks))
	for _, handMask := range handMasks {
		handValues = append(handValues, Evaluate7CardByMask(handMask|boardMask).Value)
	}
	maxIndexs := findMaxIndexs(handValues)
	if len(maxIndexs) == 1 {
		winCounts[maxIndexs[0]]++
	} else {
		for _, maxIndex := range maxIndexs {
			tieCounts[maxIndex]++
		}
	}
	count++
	return WinStat{winCounts, tieCounts, count}
}

func findMaxIndexs(arr []uint64) []int {
	maxIndexs := make([]int, 0, len(arr))
	maxValue := uint64(0)
	for _, value := range arr {
		if value > maxValue {
			maxValue = value
		}
	}
	for index, value := range arr {
		if value == maxValue {
			maxIndexs = append(maxIndexs, index)
		}
	}
	return maxIndexs
}

func (totalWinStat *WinStat) mergeWinStat(winStat WinStat) {
	totalWinStat.Count += winStat.Count
	for index, winCount := range winStat.WinCounts {
		totalWinStat.WinCounts[index] += winCount
	}
	for index, tieCount := range winStat.TieCounts {
		totalWinStat.TieCounts[index] += tieCount
	}
}
