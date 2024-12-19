package scrum_poker

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type GameState int

const (
	WaitingToStart GameState = 0
	Voting         GameState = 1
	Reveal         GameState = 2
)

type Game struct {
	ID      uuid.UUID
	Players map[uuid.UUID]*Player
	State   GameState
}

func NewGame() *Game {
	return &Game{
		ID:    uuid.New(),
		State: WaitingToStart,
	}
}

func (g *Game) newPlayerConnected(playerName string) error {
	log.Debug().Msgf("game id: %s \t- new player '%s' joined", g.ID.String(), playerName)
	newPlayer := NewPlayer(playerName)

	//FIXME: using a whole copy of the uuid object as key seems inefficient
	g.Players[newPlayer.ID] = newPlayer
	return nil
}

func (g *Game) playerDisconnected(player *Player) error {
	log.Debug().Msgf("game id: %s \t- player '%s disconnected", g.ID.String(), player.Name)
	delete(g.Players, player.ID)
	g.continueIfVoteCompleted()
	return nil
}

func (g *Game) startVoting() error {
	log.Debug().Msgf("game id: %s \t- starting vote in game", g.ID.String())
	if g.State != WaitingToStart {
		return fmt.Errorf("cannot start voting - game is not currently waiting to start")
	}

	if len(g.Players) < 1 {
		return fmt.Errorf("cannot start voting - not enough players")
	}

	g.State = Voting
	return nil
}

func (g *Game) submitVote(player *Player, vote Vote) error {
	log.Debug().Msgf("game id: %s \t- vote '%s' submitted by player name: '%s'", g.ID.String(), vote.text, player.Name)
	player.Vote = vote
	g.continueIfVoteCompleted()
	return nil
}

func (g *Game) continueIfVoteCompleted() error {
	if g.State != Voting {
		return nil
	}

	// If any players haven't voted, just return
	for _, player := range g.Players {
		if player.Vote.text == "" {
			return nil
		}
	}

	log.Debug().Msgf("game id: %s \t- all players have voted, moving on to reveal", g.ID.String())

	g.State = Reveal
	return nil
}
