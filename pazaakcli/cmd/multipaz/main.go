package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/pazaak/pazaakcli/pazaak"
	"github.com/loopfz/pazaak/pazaakcli/player"
	"github.com/sirupsen/logrus"
)

var human *player.AsyncPlayer

func main() {

	bot := flag.String("player", "", "player program path")
	quiet := flag.Bool("quiet", false, "no logs")
	flag.Parse()

	if *quiet {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	if *bot == "" {
		panic("missing bot")
	}

	aiPlayer := player.NewForkPlayer(*bot)
	human = player.NewAsyncPlayer()

	router := gin.Default()
	router.GET("/state", getState)
	router.POST("/move", doMove)

	g, err := pazaak.NewGame([]player.Player{aiPlayer, human}, "", pazaak.AutoSidedeckHandler{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}

	go func() { router.Run(":8087") }()

	g.Run()
}

func getState(c *gin.Context) {
	human.GetState(c.Writer)
}

func doMove(c *gin.Context) {

	err := human.DoMove(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(400, map[string]string{"error": err.Error()})
	}
}
