package main

import "math/rand"

type AsteroidType int

const (
	RockAsteroid AsteroidType = iota
	FuelAsteroid
)

type Asteroid struct {
	ID       int          `json:"id"`
	Position Position     `json:"position"`
	Vector   Position     `json:"vector"`
	Type     AsteroidType `json:"type"`
	Size     float64      `json:"size"`
}

func NewAsteroid(m *Map) *Asteroid {
	a := &Asteroid{
		ID:       len(m.Asteroids),
		Position: RandomPosition(m),
		Type:     AsteroidType(rand.Intn(2)),
		Size:     RandomFloat(MinAsteroidSize, MaxAsteroidSize),
	}

	m.Asteroids = append(m.Asteroids, a)
	return a
}
