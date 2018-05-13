package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/loopfz/pazaak/pazaakcli/pazaak"
	"github.com/loopfz/pazaak/pazaakcli/player"
	"github.com/sirupsen/logrus"
)

func main() {

	var playerList arrayFlags
	flag.Var(&playerList, "player", "player program path")
	statsFile := flag.String("stats", "", "stats file")
	quiet := flag.Bool("quiet", false, "no logs")
	flag.Parse()

	if *quiet {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	var pl []*player.Player
	for _, p := range playerList {
		pl = append(pl, player.NewForkPlayer(p))
	}

	g, err := pazaak.NewGame(pl, *statsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}

	g.Run()
}
