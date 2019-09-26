package holdem

import (
	"math/rand"
	"time"
)

// RandomNCardMask 随机N张card 效率较高
func RandomNCardMask(n int, dead uint64) []uint64 {
	cardMasks := make([]uint64, n)

	shuffleCards := ShuffleCardMasks()
	i := 0
	for _, card := range shuffleCards {
		if card&dead == 0 {
			cardMasks[i] = card
			i++
		}
		if i >= n {
			break
		}
	}
	return cardMasks
}

// RandomNCard 随机N张card 效率较低
func RandomNCard(n int, dead uint64) []string {
	cards := make([]string, n)

	shuffleCards := ShuffleCards()
	i := 0
	for _, card := range shuffleCards {
		if ParseHand([]string{card})&dead == 0 {
			cards[i] = card
			i++
		}
		if i >= n {
			break
		}
	}

	return cards
}

// ShuffleCardMasks 打乱Card顺序
func ShuffleCardMasks() []uint64 {
	r := myRand()
	ret := make([]uint64, len(cardMasksTable))
	perm := r.Perm(len(cardMasksTable))
	for i, randIndex := range perm {
		ret[i] = cardMasksTable[randIndex]
	}
	return ret
}

// ShuffleCards 打乱Card顺序
func ShuffleCards() []string {
	r := myRand()
	ret := make([]string, len(cardTable))
	perm := r.Perm(len(cardTable))
	for i, randIndex := range perm {
		ret[i] = cardTable[randIndex]
	}
	return ret
}

func myRand() *rand.Rand {
	s := rand.NewSource(time.Now().UnixNano())
	return rand.New(s)
}
