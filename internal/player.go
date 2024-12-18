package main

import (
	"github.com/google/uuid"
)

type Player struct {
	Name string
	ID   uuid.UUID
	Vote Vote
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		ID:   uuid.New(),
	}
}
