package main

import (
	"fmt"

	"github.com/lixw1994/pokerevaluator/holdem"
)

func main() {
	fmt.Println("holdem evaluator")
	cards := []string{"as", "ad"}
	cardMask := holdem.ParseHand(cards)
	// boards := []string{"ah", "Qs", "Js", "Ts", "8d"}
	boards := []string{"Ts", "Qs", "2d", "6c"}
	// boards := []string{"Ts", "Qs", "2d"}
	// boards := []string{}
	boardMask := holdem.ParseHand(boards)
	fmt.Println(cardMask, boardMask, holdem.StringifyHand(cardMask), holdem.StringifyHand(boardMask))

	holdemGame := holdem.NewGame(5)
	fmt.Println(holdemGame.Players, holdemGame.Cards, holdemGame.Boards)
}
