package scrum_poker

import "github.com/google/uuid"

type GameController struct {
	Games map[uuid.UUID]*Game
}
