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

	target := 20
	fallbackTarget := 19
	if len(game.CurrentPlayer.Hand) == 0 {
		fallbackTarget = 17
	}

	if game.Opponent.Stand {
		if game.Opponent.BoardValue == 20 {
			target = 20
			fallbackTarget = 20
		} else {
			target = game.Opponent.BoardValue + 1
			fallbackTarget = game.Opponent.BoardValue
		}
	}

	if matchTarget(game.CurrentPlayer.BoardValue, target) {
		move.Stand = true
	} else {

		handCard, flip := findCard(game, target)
		if handCard != "" {
			playCard(handCard, flip, move)
		} else {
			if matchTarget(game.CurrentPlayer.BoardValue, fallbackTarget) {
				move.Stand = true
			} else {
				handCard, flip = findCard(game, fallbackTarget)
				if handCard != "" {
					playCard(handCard, flip, move)
				}
			}
		}
	}

	if move.HandCard == "" && game.CurrentPlayer.BoardValue > 20 && game.Opponent.RoundWins == pazaak.MAX_ROUND_WINS-1 {
		// last resort, try to survive
		tmpVal := 0
		handCard := ""
		flipCard := false
		for _, c := range game.CurrentPlayer.Hand {
			val := game.CurrentPlayer.BoardValue + c.Value
			if c.Flip {
				val = game.CurrentPlayer.BoardValue - c.Value
			}
			if val < 20 && val >= tmpVal {
				tmpVal = val
				handCard = c.Identifier
				flipCard = c.Flip
			}
		}
		if handCard != "" {
			move.HandCard = handCard
			move.FlipCard = flipCard
			if !game.Opponent.Stand {
				if tmpVal >= 17 && len(game.CurrentPlayer.Hand)-1 == 0 {
					move.Stand = true
				}
			}
		}
	}

	j, err := json.Marshal(move)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}

func playCard(handCard string, flip bool, move *pazaak.PazaakMove) {
	move.Stand = true
	move.HandCard = handCard
	move.FlipCard = flip
}

func findCard(game *pazaak.PazaakGame, target int) (string, bool) {

	for _, c := range game.CurrentPlayer.Hand {
		if matchTarget(game.CurrentPlayer.BoardValue+c.Value, target) {
			return c.Identifier, false
		}
		if c.Flip && matchTarget(game.CurrentPlayer.BoardValue-c.Value, target) {
			return c.Identifier, true
		}
	}
	return "", false
}

func matchTarget(val int, target int) bool {
	return val >= target && val <= 20
}
