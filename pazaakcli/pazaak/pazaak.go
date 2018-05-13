package pazaak

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/loopfz/pazaak/pazaakcli/player"
	"github.com/sirupsen/logrus"
)

const (
	SIDEDECK_SIZE   = 10
	HAND_SIZE       = 4
	MAX_BOARD_SIZE  = 9
	MAX_BOARD_VALUE = 20
	MAX_ROUND_WINS  = 3

	AUTO_SIDEDECK = "auto"
)

type PazaakGame struct {
	StatsFile     string          `json:"-"`
	Players       []*PazaakPlayer `json:"-"`
	CurrentPlayer PazaakPlayer    `json:"current_player"`
	Opponent      PazaakPlayer    `json:"opponent"`
	Winner        PazaakPlayer    `json:"winner"`
}

type PazaakPlayer struct {
	*player.Player

	// Not reset between rounds
	SideDeck    []*PazaakCard `json:"-"`
	InitialHand []*PazaakCard `json:"initial_hand"`
	Hand        []*PazaakCard `json:"hand"`
	RoundWins   int           `json:"round_wins"`
	Winner      bool          `json:"winner"`

	// Reset every round
	Deck       []*PazaakCard `json:"-"`
	Board      []*PazaakCard `json:"board"`
	BoardValue int           `json:"board_value"`
	Stand      bool          `json:"stand"`
}

type PazaakCard struct {
	Identifier string `json:"identifier"`
	Value      int    `json:"value"`
	Flip       bool   `json:"flip,omitempty"`
	Special    bool   `json:"special,omitempty"`
}

type PazaakMove struct {
	HandCard string `json:"hand_card"`
	FlipCard bool   `json:"flip_card"`
	Stand    bool   `json:"stand"`
}

type Stats struct {
	Score map[string]int `json:"score"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (m *PazaakMove) Valid() error {
	return nil
}

func NewGame(pl []*player.Player, statsFile string) (*PazaakGame, error) {

	g := &PazaakGame{StatsFile: statsFile}

	for _, p := range pl {
		p.Number = uint(len(g.Players) + 1)
		g.Players = append(g.Players, &PazaakPlayer{Player: p})
	}

	if len(g.Players) != 2 {
		return nil, errors.New("Player count should be 2")
	}

	err := g.InitPlayerSideDecks()
	if err != nil {
		return nil, err
	}

	g.InitPlayerHands()
	g.InitFirstPlayer(rand.Intn(len(g.Players)))

	return g, nil
}

var knownCards = map[string]PazaakCard{
	"+1":  {Value: 1, Identifier: "+1"},
	"+2":  {Value: 2, Identifier: "+2"},
	"+3":  {Value: 3, Identifier: "+3"},
	"+4":  {Value: 4, Identifier: "+4"},
	"+5":  {Value: 5, Identifier: "+5"},
	"+6":  {Value: 6, Identifier: "+6"},
	"-1":  {Value: -1, Identifier: "-1"},
	"-2":  {Value: -2, Identifier: "-2"},
	"-3":  {Value: -3, Identifier: "-3"},
	"-4":  {Value: -4, Identifier: "-4"},
	"-5":  {Value: -5, Identifier: "-5"},
	"-6":  {Value: -6, Identifier: "-6"},
	"+-1": {Value: 1, Identifier: "+-1", Flip: true},
	"+-2": {Value: 2, Identifier: "+-2", Flip: true},
	"+-3": {Value: 3, Identifier: "+-3", Flip: true},
	"+-4": {Value: 4, Identifier: "+-4", Flip: true},
	"+-5": {Value: 5, Identifier: "+-5", Flip: true},
	"+-6": {Value: 6, Identifier: "+-6", Flip: true},
}

func NewPazaakCard(ident string) (*PazaakCard, error) {
	ident = strings.TrimSpace(ident)
	c, ok := knownCards[ident]
	if !ok {
		return nil, fmt.Errorf("Unknown card '%s'", ident)
	}
	return &c, nil
}

func NewPazaakDeck() []*PazaakCard {
	ret := []*PazaakCard{
		{Value: 1, Identifier: "1"},
		{Value: 1, Identifier: "1"},
		{Value: 1, Identifier: "1"},
		{Value: 1, Identifier: "1"},
		{Value: 2, Identifier: "2"},
		{Value: 2, Identifier: "2"},
		{Value: 2, Identifier: "2"},
		{Value: 2, Identifier: "2"},
		{Value: 3, Identifier: "3"},
		{Value: 3, Identifier: "3"},
		{Value: 3, Identifier: "3"},
		{Value: 3, Identifier: "3"},
		{Value: 4, Identifier: "4"},
		{Value: 4, Identifier: "4"},
		{Value: 4, Identifier: "4"},
		{Value: 4, Identifier: "4"},
		{Value: 5, Identifier: "5"},
		{Value: 5, Identifier: "5"},
		{Value: 5, Identifier: "5"},
		{Value: 5, Identifier: "5"},
		{Value: 6, Identifier: "6"},
		{Value: 6, Identifier: "6"},
		{Value: 6, Identifier: "6"},
		{Value: 6, Identifier: "6"},
		{Value: 7, Identifier: "7"},
		{Value: 7, Identifier: "7"},
		{Value: 7, Identifier: "7"},
		{Value: 7, Identifier: "7"},
		{Value: 8, Identifier: "8"},
		{Value: 8, Identifier: "8"},
		{Value: 8, Identifier: "8"},
		{Value: 8, Identifier: "8"},
		{Value: 9, Identifier: "9"},
		{Value: 9, Identifier: "9"},
		{Value: 9, Identifier: "9"},
		{Value: 9, Identifier: "9"},
		{Value: 10, Identifier: "10"},
		{Value: 10, Identifier: "10"},
		{Value: 10, Identifier: "10"},
		{Value: 10, Identifier: "10"},
	}

	for i := len(ret) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret
}

func (g *PazaakGame) NewMove() player.PlayerMove {
	return &PazaakMove{}
}

func buildRandomSideDeck() []string {
	var ret []string
	keys := []string{}
	for k, _ := range knownCards {
		keys = append(keys, k)
	}
	for i := 0; i < SIDEDECK_SIZE; i++ {
		ret = append(ret, keys[rand.Intn(len(keys))])
	}
	return ret
}

func (g *PazaakGame) InitPlayerSideDecks() error {
	reader := bufio.NewReader(os.Stdin)
	for _, p := range g.Players {
		s, _ := reader.ReadString('\n')
		s = strings.TrimSpace(s)
		var cards []string
		if s == AUTO_SIDEDECK {
			cards = buildRandomSideDeck()
		} else {
			cards = strings.Split(s, ",")
		}
		if len(cards) != SIDEDECK_SIZE {
			return fmt.Errorf("%v: invalid side deck, expected %d elements", cards, SIDEDECK_SIZE)
		}
		for _, c := range cards {
			newCard, err := NewPazaakCard(c)
			if err != nil {
				return err
			}
			p.SideDeck = append(p.SideDeck, newCard)
		}
		for i := len(p.SideDeck) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			p.SideDeck[i], p.SideDeck[j] = p.SideDeck[j], p.SideDeck[i]
		}
	}
	return nil
}

func (g *PazaakGame) InitPlayerHands() {
	for _, p := range g.Players {
		if HAND_SIZE > len(p.SideDeck) {
			panic("handsize > side deck")
		}
		for i := 0; i < HAND_SIZE; i++ {
			p.Hand = append(p.Hand, p.SideDeck[i])
		}
		p.InitialHand = p.Hand
	}
}

// factorize?
func (g *PazaakGame) InitFirstPlayer(i int) {
	pl := g.Players[i:]
	for j := 0; j < i; j++ {
		pl = append(pl, g.Players[j])
	}
	g.Players = pl
}

func (g *PazaakGame) Run() {

GAMELOOP:
	for {
		logrus.Infof("--- NEW ROUND ---\nPlayer %s is first player!\n", g.Players[0])

		// reset round stuff
		for _, p := range g.Players {
			p.Deck = NewPazaakDeck()
			p.Board = nil
			p.BoardValue = 0
			p.Stand = false
		}
	ROUNDLOOP:
		for {
			for i, p := range g.Players {
				opponentIdx := (i + 1) % 2
				g.Opponent = *(g.Players[opponentIdx])
				g.Opponent.Hand = nil
				if !p.Stand {
					p.DrawCard()
					g.CurrentPlayer = *p
					move, err := p.GetMove(g)
					if err != nil {
						logrus.Infof("Player %s move error: %s\n", p, err)
					} else {
						pzkMove, ok := move.(*PazaakMove)
						if !ok {
							panic("player returned a non-pazaak move")
						}
						if pzkMove.HandCard != "" {
							p.PlayHandCard(pzkMove.HandCard, pzkMove.FlipCard)
						}
						if pzkMove.Stand {
							logrus.Infof("Player %s STANDS\n", p)
							p.Stand = true
						}
					}
				}
				if !p.Stand && p.BoardValue == MAX_BOARD_VALUE {
					logrus.Infof("Player %s AUTO STANDS\n", p)
					p.Stand = true
				}
				if !p.Stand {
					logrus.Infof("Player %s CONTINUES\n", p)
				}
				if p.BoardValue > MAX_BOARD_VALUE {
					logrus.Infof("Player %s busted, opponent scores\n", p)
					g.Players[opponentIdx].RoundWins++
					g.InitFirstPlayer(opponentIdx)
					break ROUNDLOOP
				}
				if p.Stand && g.Opponent.Stand {
					if p.BoardValue == g.Opponent.BoardValue {
						break ROUNDLOOP
					}
					winnerIdx := opponentIdx
					if p.BoardValue > g.Opponent.BoardValue {
						winnerIdx = i
					}
					g.Players[winnerIdx].RoundWins++
					g.InitFirstPlayer(winnerIdx)
					break ROUNDLOOP
				}
			}
		}
		for _, p := range g.Players {
			if p.RoundWins >= MAX_ROUND_WINS {
				p.Winner = true
				g.CurrentPlayer = PazaakPlayer{}
				g.Opponent = PazaakPlayer{}
				g.Winner = *p
				break GAMELOOP
			}
		}
	}

	g.WriteStats()

	logrus.Infof("Player %s WINS", g.Winner)

	os.Exit(int(g.Winner.Number))
}

func (g *PazaakGame) ReadStats() (*Stats, error) {
	ret := &Stats{Score: map[string]int{}}

	if g.StatsFile == "" {
		return ret, nil
	}

	b, err := ioutil.ReadFile(g.StatsFile)
	if err != nil {
		return ret, nil
	}
	err = json.Unmarshal(b, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (g *PazaakGame) WriteStats() {
	if g.StatsFile == "" {
		return
	}
	stats, err := g.ReadStats()
	if err != nil {
		panic(err)
	}
	for _, p := range g.Players {
		delta := 1
		if !p.Winner {
			delta = -1
		}
		for _, c := range p.InitialHand {
			stats.Score[c.Identifier] = stats.Score[c.Identifier] + delta
		}
	}
	j, err := json.Marshal(stats)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(g.StatsFile, j, 0644)
}

func (p *PazaakPlayer) PlayHandCard(ident string, flipCard bool) {
	var playedCard *PazaakCard
	for i, c := range p.Hand {
		if c.Identifier == ident {
			playedCard = p.Hand[i]
			p.Hand[i] = p.Hand[len(p.Hand)-1]
			p.Hand = p.Hand[:len(p.Hand)-1]
			break
		}
	}
	if playedCard != nil {
		if flipCard {
			if !playedCard.Flip {
				logrus.Infof("Player %s requested to play card %s flipped, but card is not flippable\n", p, ident)
				return
			} else {
				playedCard.Value = -playedCard.Value
			}
		}
		logrus.Infof("Player %s plays hand card %s\n", p, ident)
		p.AddBoardCard(playedCard)
	} else if ident != "" {
		logrus.Infof("Player %s requested to play card %s, but does not have it in hand\n", p, ident)
	}
}

func (p *PazaakPlayer) DrawCard() {

	if len(p.Deck) == 0 {
		panic("trying to draw from empty deck")
	}
	newCard := p.Deck[0]
	p.Deck = p.Deck[1:]

	logrus.Infof("Player %s draws a %s\n", p, newCard.Identifier)

	p.AddBoardCard(newCard)
}

func (p *PazaakPlayer) AddBoardCard(newCard *PazaakCard) {

	p.Board = append(p.Board, newCard)
	p.BoardValue += newCard.Value
}

func (p PazaakPlayer) String() string {
	hand := []string{}
	for _, c := range p.Hand {
		hand = append(hand, c.Identifier)
	}
	return fmt.Sprintf("%s [%d] {%v}", p.Player, p.BoardValue, hand)
}
