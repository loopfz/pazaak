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

	if game.CurrentPlayer.BoardValue >= 15 || (game.Opponent.Stand && game.CurrentPlayer.BoardValue >= game.Opponent.BoardValue) {
		move.Stand = true
	}

	j, err := json.Marshal(move)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
