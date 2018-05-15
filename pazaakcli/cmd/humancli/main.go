package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/loopfz/pazaak/pazaakcli/pazaak"
)

func main() {

	resp, err := http.Get("http://localhost:8087/state")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	resp.Body.Close()

	game := pazaak.PazaakGame{}

	err = json.Unmarshal(body, &game)
	if err != nil {
		panic(err)
	}

	var handIdent []string
	for _, c := range game.CurrentPlayer.Hand {
		handIdent = append(handIdent, c.Identifier)
	}

	fmt.Println("----------")
	fmt.Println("OPPONENT")
	fmt.Println("----------")
	fmt.Printf("Round wins: %d\n", game.Opponent.RoundWins)
	fmt.Printf("Board value: %d\n", game.Opponent.BoardValue)
	fmt.Printf("Stands: %v\n", game.Opponent.Stand)

	fmt.Println("----------")
	fmt.Println("YOU")
	fmt.Println("----------")
	fmt.Printf("Round wins: %d\n", game.CurrentPlayer.RoundWins)
	fmt.Printf("Board value: %d\n", game.CurrentPlayer.BoardValue)
	fmt.Printf("Hand: %s\n", strings.Join(handIdent, ", "))

	move := &pazaak.PazaakMove{}

	move.HandCard = prompter.Prompt("Play hand card?", "")
	if strings.HasPrefix(move.HandCard, "+-") {
		move.FlipCard = prompter.YN("Flip card (use negative value) ?", false)
	}
	move.Stand = prompter.YN("Stand?", false)

	movestr, err := json.Marshal(move)
	if err != nil {
		panic(err)
	}

	_, err = http.Post("http://localhost:8087/move", "application/json", bytes.NewBuffer(movestr))
	if err != nil {
		panic(err)
	}

}
