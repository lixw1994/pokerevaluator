package holdem

import "fmt"

// Game 牌局
type Game struct {
	Cards   [][]string
	Boards  []string
	Players []Player
}

// Player 玩家
type Player struct {
	Cards        []string
	HandValue    uint64
	MaxHand      []string
	HandType     uint64
	HandTypeDesc string
}

// NewGame 模拟牌局
func NewGame(playerNum int) (game Game) {
	game.Cards = make([][]string, 0, playerNum)
	game.Boards = make([]string, numberOfBoard)
	game.Players = make([]Player, 0, playerNum)
	n := playerNum*numberOfHand + numberOfBoard
	randomCards := RandomNCard(n, 0)
	fmt.Println(randomCards)

	game.Boards = randomCards[n-5:]
	for index := 0; index < playerNum*numberOfHand; index += numberOfHand {
		cards := randomCards[index : index+numberOfHand]
		game.Cards = append(game.Cards, cards)
		game.Players = append(game.Players, newPlayer(cards, game.Boards))
	}
	// game.Boards = randomCards[index : index+numberOfBoard]
	return
}

func newPlayer(cards []string, boards []string) Player {
	res := Evaluate7CardByMask(ParseHand(cards) | ParseHand(boards))
	handType := HandType(res.Value)
	handTypeDesc := HandTypeDesc(handType)
	return Player{
		Cards:        cards,
		HandValue:    res.Value,
		MaxHand:      res.MaxHands,
		HandType:     handType,
		HandTypeDesc: handTypeDesc,
	}
}
