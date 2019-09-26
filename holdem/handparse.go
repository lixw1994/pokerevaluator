package holdem

import "errors"

// ParseHand 解析字符串 handMask
func ParseHand(cards []string) (cardMask uint64) {
	for _, card := range cards {
		var rank, suit uint
		switch card[0] {
		case '2':
			rank = rank2
		case '3':
			rank = rank3
		case '4':
			rank = rank4
		case '5':
			rank = rank5
		case '6':
			rank = rank6
		case '7':
			rank = rank7
		case '8':
			rank = rank8
		case '9':
			rank = rank9
		case 'T':
			fallthrough
		case 't':
			rank = rankTen
		case 'J':
			fallthrough
		case 'j':
			rank = rankJack
		case 'Q':
			fallthrough
		case 'q':
			rank = rankQueen
		case 'K':
			fallthrough
		case 'k':
			rank = rankKing
		case 'A':
			fallthrough
		case 'a':
			rank = rankAce
		default:
			panic(errors.New("rank param error"))
		}
		switch card[1] {
		case 'H':
			fallthrough
		case 'h':
			suit = heart
		case 'D':
			fallthrough
		case 'd':
			suit = diamond
		case 'C':
			fallthrough
		case 'c':
			suit = club
		case 'S':
			fallthrough
		case 's':
			suit = spade
		default:
			panic(errors.New("suit param error"))
		}

		cardMask = cardMask | (uint64(1) << (rank + (suit * 13)))
	}

	return cardMask
}

// StringifyHand 生产字符串 cards
func StringifyHand(cardMask uint64) (cards []string) {
	for i := 0; cardMask != 0; i++ {
		if (cardMask & 0x1) != 0 {
			cards = append(cards, cardTable[i])
		}
		cardMask = cardMask >> 1
	}
	return
}

// HandType 牌型
func HandType(handValue uint64) uint64 {
	return handValue >> handtypeShift
}

// HandTypeDesc 牌型描述
func HandTypeDesc(handType uint64) string {
	var handTypeDesc string
	switch handType {
	case HighCard:
		handTypeDesc = "HighCard"
	case Pair:
		handTypeDesc = "Pair"
	case TwoPair:
		handTypeDesc = "TwoPair"
	case Trips:
		handTypeDesc = "Trips"
	case Straight:
		handTypeDesc = "Straight"
	case Flush:
		handTypeDesc = "Flush"
	case FullHouse:
		handTypeDesc = "FullHouse"
	case FourOfAKind:
		handTypeDesc = "FourOfAKind"
	case StraightFlush:
		handTypeDesc = "StraightFlush"
	default:
		handTypeDesc = ""
	}
	return handTypeDesc
}

// PocketType 获取pocket 169类型 one of Pocket Hand 169 Enum
func PocketType(pocketMask uint64) int {
	for typeIndex, pocketMasks := range pocket169Table {
		for _, pm := range pocketMasks {
			if pm == pocketMask {
				return typeIndex
			}
		}
	}
	return -1
}
