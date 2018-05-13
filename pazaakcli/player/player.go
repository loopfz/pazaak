package player

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	PLAYER_TIMEOUT = 5 * time.Second
)

type Player struct {
	Program string `json:"-"`
	Number  uint   `json:"number"`
	// Data     *json.RawMessage                              `json:"data"` // TODO
	executor func(*Player, GameEngine) (PlayerMove, error) `json:"-"`
}

type GameEngine interface {
	NewMove() PlayerMove
}

type PlayerMove interface {
	Valid() error
}

func NewForkPlayer(bin string) *Player {
	return &Player{
		Program:  bin,
		executor: ForkMove,
	}
}

func (p *Player) String() string {
	return fmt.Sprintf("%d (%s)", p.Number, p.Program)
}

func (p *Player) GetMove(g GameEngine) (PlayerMove, error) {
	if p.executor == nil {
		return nil, errors.New("No executor set")
	}
	return p.executor(p, g)
}

func ForkMove(p *Player, g GameEngine) (PlayerMove, error) {

	cmd := exec.Command(p.Program)

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	j, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}
	in.Write(j)
	cmd.Stdin = in
	cmd.Stdout = out

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(PLAYER_TIMEOUT):
		if err := cmd.Process.Kill(); err != nil {
			fmt.Fprintf(os.Stderr, "Player %s (%d): Failed to kill subprocess after timeout: %s", p.Number, p.Program, err)
		}
		<-done
		return nil, errors.New("Timeout")
	case err := <-done:
		if err != nil {
			return nil, err
		}
	}

	move := g.NewMove()
	err = json.Unmarshal(out.Bytes(), move)
	if err != nil {
		return nil, err
	}

	err = move.Valid()
	if err != nil {
		return nil, err
	}

	return move, nil
}
