package player

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	PLAYER_TIMEOUT = 5 * time.Second
)

type Player interface {
	GetMove(GameEngine) (PlayerMove, error)
}

type AIPlayer struct {
	Program string `json:"-"`
	// Data     *json.RawMessage                              `json:"data"` // TODO
	executor func(*AIPlayer, GameEngine) (PlayerMove, error) `json:"-"`
}

type GameEngine interface {
	NewMove() PlayerMove
}

type PlayerMove interface {
	Valid() error
}

func NewForkPlayer(bin string) Player {
	return &AIPlayer{
		Program:  bin,
		executor: ForkMove,
	}
}

func (p *AIPlayer) String() string {
	return p.Program
}

func (p *AIPlayer) GetMove(g GameEngine) (PlayerMove, error) {
	if p.executor == nil {
		return nil, errors.New("No executor set")
	}
	return p.executor(p, g)
}

func ForkMove(p *AIPlayer, g GameEngine) (PlayerMove, error) {

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
			fmt.Fprintf(os.Stderr, "Player %s: Failed to kill subprocess after timeout: %s", p.Program, err)
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

type AsyncPlayer struct {
	mut       sync.Mutex
	gameState []byte
	deadline  time.Time
	g         GameEngine
	respChan  chan PlayerMove
}

func (ap *AsyncPlayer) GetMove(g GameEngine) (PlayerMove, error) {

	ap.mut.Lock()
	rawG, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}
	ap.gameState = rawG
	ap.g = g
	deadline := time.Now().Add(1 * time.Minute)
	ap.deadline = deadline
	ch := ap.respChan
	ap.mut.Unlock()
	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(1 * time.Minute):
	}

	return nil, errors.New("human timeout")
}

func (ap *AsyncPlayer) GetState(w http.ResponseWriter) {
	ap.mut.Lock()
	w.Write(ap.gameState)
	ap.mut.Unlock()
}

func (ap *AsyncPlayer) DoMove(r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	ap.mut.Lock()
	defer ap.mut.Unlock()
	move := ap.g.NewMove()
	err = json.Unmarshal(body, move)
	if err != nil {
		return err
	}
	err = move.Valid()
	if err != nil {
		return err
	}
	select {
	case ap.respChan <- move:
		return nil
	default:
	}

	return errors.New("too late to play!")
}
