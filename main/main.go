package main

import (
	"fmt"

	"github.com/lixw1994/pokerevaluator/holdem"
)

func main() {
	fmt.Println("holdem evaluator")
	cards := []string{"5s", "8s"}
	cardMask := holdem.ParseHand(cards)
	boards := []string{"4s", "6s", "Qs", "7s", "Tc"}
	// boards := []string{"Ts", "Qs", "2d", "6c"}
	// boards := []string{"Ts", "Qs", "2d"}
	// boards := []string{}
	boardMask := holdem.ParseHand(boards)
	res := holdem.Evaluate7CardByMask(cardMask | boardMask)
	fmt.Println(res, holdem.HandTypeDesc(holdem.HandType(res.Value)))
	fmt.Println(cardMask, boardMask, holdem.StringifyHand(cardMask), holdem.StringifyHand(boardMask))

	holdemGame := holdem.NewGame(5)
	fmt.Println(holdemGame.Players, holdemGame.Cards, holdemGame.Boards)
}
