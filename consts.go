package main

import "math/rand"

const (
	Radius          = 15000
	MaxAsteroidSize = 50
	MinAsteroidSize = MaxAsteroidSize / 2
	AsteroidCount   = 500
	WormholeCount   = 25
	ShipMaxHealth   = 100
	ShipStartFuel   = 100
	PlayerStartFuel = 1000
	PlayerStartRock = 1000
)

func ShipRockPrice(t ShipType) int {
	return 100
}

func RandomFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
