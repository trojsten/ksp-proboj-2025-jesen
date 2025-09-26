package main

import (
	"math"
	"math/rand"
)

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

func NewAsteroidFromShip(m *Map, ship *Ship, asteroidType AsteroidType) *Asteroid {
	var materialAmount float64
	if asteroidType == FuelAsteroid {
		materialAmount = ship.Fuel
	} else {
		materialAmount = float64(ship.Rock)
	}

	size := math.Sqrt(materialAmount / MaterialToSurfaceRatio / math.Pi)

	a := &Asteroid{
		ID:       len(m.Asteroids),
		Position: RandomOffsetPosition(ship.Position, AsteroidSpawnOffset),
		Type:     asteroidType,
		Size:     size,
	}

	m.Asteroids = append(m.Asteroids, a)
	return a
}
