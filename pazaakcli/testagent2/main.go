package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/loopfz/pazaak/pazaakcli/pazaak"
)

func main() {

	decoder := json.NewDecoder(os.Stdin)

	game := &pazaak.PazaakGame{}
	err := decoder.Decode(game)
	if err != nil {
		panic(err)
	}

	move := &pazaak.PazaakMove{}

	// aim for 20 at the cost of cards
	// but keep a natural 19
	commitCardThreshold := 20
	standThreshold := 19

	// no cards in hand, or last card may get used just to survive
	// keep a natural 17 to limit busts
	if len(game.CurrentPlayer.Hand) == 0 ||
		(len(game.CurrentPlayer.Hand) == 1 && game.CurrentPlayer.BoardValue > 20) {
		standThreshold = 17
	}

	// if opp stands, aim for value+1 with cards
	// but keep anything that ties
	if game.Opponent.Stand {
		commitCardThreshold = game.Opponent.BoardValue + 1
		standThreshold = game.Opponent.BoardValue
	}

	// DANGER, commit cards to reach our lowest acceptable value
	// if currently over 20, spend anything to survive
	if game.Opponent.RoundWins == 2 {
		commitCardThreshold = standThreshold
		if game.CurrentPlayer.BoardValue > 20 {
			commitCardThreshold = 0
		}
	}

	if commitCardThreshold > 20 {
		commitCardThreshold = 20
	}

	valueModifier := 0

	// if we're not currently naturally over our ideal target, find the best hand card
	if !legalValueOver(game.CurrentPlayer.BoardValue, commitCardThreshold) {
		handCard := ""
		flipCard := false
		val := 0
		modif := 0
		for _, c := range game.CurrentPlayer.Hand {
			tmpVal := game.CurrentPlayer.BoardValue + c.Value
			tmpFlip := false
			if tmpVal > 20 && c.Flip {
				tmpVal = game.CurrentPlayer.BoardValue - c.Value
				tmpFlip = true
			}
			if tmpVal > val && legalValueOver(tmpVal, commitCardThreshold) {
				handCard = c.Identifier
				flipCard = tmpFlip
				val = tmpVal
				modif = c.Value
				if tmpFlip {
					modif = -c.Value
				}
			}
		}
		if handCard != "" {
			move.HandCard = handCard
			move.FlipCard = flipCard
			valueModifier = modif
		}
	}

	// if (optionally with the card we're playing) we're over our minimum value to stand, do it
	if legalValueOver(game.CurrentPlayer.BoardValue+valueModifier, standThreshold) {
		move.Stand = true
	}

	j, err := json.Marshal(move)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}

func legalValueOver(val int, target int) bool {
	return val >= target && val <= 20
}
