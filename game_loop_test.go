package main

import (
	"encoding/json"
	"testing"
)

func TestGameStateFor(t *testing.T) {
	m := &Map{
		Radius: 100,
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Position: Position{X: 10, Y: 20}},
		},
		Asteroids: []*Asteroid{
			{ID: 0, Position: Position{X: 30, Y: 40}},
		},
		Wormholes: []*Wormhole{
			{ID: 0, Position: Position{X: 50, Y: 60}},
		},
		Players: []*Player{
			{ID: 0, Name: "TestPlayer"},
		},
		Round: 5,
	}
	p := &Player{ID: 0, Name: "TestPlayer"}

	state := GameStateFor(m, p)

	var decodedMap Map
	err := json.Unmarshal([]byte(state), &decodedMap)
	if err != nil {
		t.Errorf("GameStateFor() returned invalid JSON: %v", err)
	}

	if decodedMap.Radius != m.Radius {
		t.Errorf("GameStateFor() Radius = %v, want %v", decodedMap.Radius, m.Radius)
	}
	if len(decodedMap.Ships) != len(m.Ships) {
		t.Errorf("GameStateFor() Ships count = %v, want %v", len(decodedMap.Ships), len(m.Ships))
	}
	if len(decodedMap.Asteroids) != len(m.Asteroids) {
		t.Errorf("GameStateFor() Asteroids count = %v, want %v", len(decodedMap.Asteroids), len(m.Asteroids))
	}
	if len(decodedMap.Wormholes) != len(m.Wormholes) {
		t.Errorf("GameStateFor() Wormholes count = %v, want %v", len(decodedMap.Wormholes), len(m.Wormholes))
	}
	if len(decodedMap.Players) != len(m.Players) {
		t.Errorf("GameStateFor() Players count = %v, want %v", len(decodedMap.Players), len(m.Players))
	}
	if decodedMap.Round != m.Round {
		t.Errorf("GameStateFor() Round = %v, want %v", decodedMap.Round, m.Round)
	}
}

func TestHandleShipMining_DrillShip(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Type: RockAsteroid, Size: 10, Position: Position{X: 5, Y: 0}},
		},
	}
	ship := &Ship{
		Type:     DrillShip,
		Position: Position{X: 0, Y: 0},
		Rock:     0,
		Fuel:     50,
	}

	HandleShipMining(m, ship)

	if ship.Rock == 0 {
		t.Errorf("HandleShipMining() ship Rock should increase, got %v", ship.Rock)
	}
}

func TestHandleShipMining_SuckerShip(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Type: FuelAsteroid, Size: 10, Position: Position{X: 5, Y: 0}},
		},
	}
	ship := &Ship{
		Type:     SuckerShip,
		Position: Position{X: 0, Y: 0},
		Rock:     0,
		Fuel:     50,
	}

	HandleShipMining(m, ship)

	if ship.Fuel == 50 {
		t.Errorf("HandleShipMining() ship Fuel should increase, got %v", ship.Fuel)
	}
}

func TestHandleShipMining_NonMiningShip(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Type: RockAsteroid, Size: 10, Position: Position{X: 20, Y: 0}},
		},
	}
	ship := &Ship{
		Type:     BattleShip,
		Position: Position{X: 0, Y: 0},
		Rock:     0,
		Fuel:     50,
	}

	HandleShipMining(m, ship)

	if ship.Rock != 0 {
		t.Errorf("HandleShipMining() non-mining ship should not mine, got Rock = %v", ship.Rock)
	}
	if ship.Fuel != 50 {
		t.Errorf("HandleShipMining() non-mining ship should not mine, got Fuel = %v", ship.Fuel)
	}
}

func TestHandleShipMining_NoAsteroids(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
	}
	ship := &Ship{
		Type:     DrillShip,
		Position: Position{X: 0, Y: 0},
		Rock:     0,
		Fuel:     50,
	}

	HandleShipMining(m, ship)

	if ship.Rock != 0 {
		t.Errorf("HandleShipMining() should not mine when no asteroids, got Rock = %v", ship.Rock)
	}
}

func TestHandleShipMining_TooFar(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Type: RockAsteroid, Size: 10, Position: Position{X: 50, Y: 0}},
		},
	}
	ship := &Ship{
		Type:     DrillShip,
		Position: Position{X: 0, Y: 0},
		Rock:     0,
		Fuel:     50,
	}

	HandleShipMining(m, ship)

	if ship.Rock != 0 {
		t.Errorf("HandleShipMining() should not mine when too far, got Rock = %v", ship.Rock)
	}
}

func TestHandleShipConquering_OwnerShip(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Size: 10, OwnerID: 0, OwnedSurface: 50, Position: Position{X: 5, Y: 0}},
		},
	}
	ship := &Ship{
		PlayerID: 0,
		Position: Position{X: 0, Y: 0},
	}

	totalSurface := 10 * 10 * 3.14159265359
	originalSurface := m.Asteroids[0].OwnedSurface

	HandleShipConquering(m, ship)

	if m.Asteroids[0].OwnedSurface <= originalSurface {
		t.Errorf("HandleShipConquering() OwnedSurface should increase for owner, got %v", m.Asteroids[0].OwnedSurface)
	}
	if m.Asteroids[0].OwnedSurface > totalSurface {
		t.Errorf("HandleShipConquering() OwnedSurface should not exceed total surface, got %v", m.Asteroids[0].OwnedSurface)
	}
}

func TestHandleShipConquering_NonOwnerShip(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Size: 10, OwnerID: 1, OwnedSurface: 50, Position: Position{X: 5, Y: 0}},
		},
	}
	ship := &Ship{
		PlayerID: 0,
		Position: Position{X: 0, Y: 0},
	}

	originalSurface := m.Asteroids[0].OwnedSurface

	HandleShipConquering(m, ship)

	if m.Asteroids[0].OwnedSurface >= originalSurface {
		t.Errorf("HandleShipConquering() OwnedSurface should decrease for non-owner, got %v", m.Asteroids[0].OwnedSurface)
	}
}

func TestHandleShipConquering_OwnerChange(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Size: 10, OwnerID: 1, OwnedSurface: 1, Position: Position{X: 5, Y: 0}},
		},
	}
	ship := &Ship{
		PlayerID: 0,
		Position: Position{X: 0, Y: 0},
	}

	HandleShipConquering(m, ship)

	if m.Asteroids[0].OwnerID != 0 {
		t.Errorf("HandleShipConquering() OwnerID should change to 0, got %v", m.Asteroids[0].OwnerID)
	}
}

func TestHandleShipConquering_NoAsteroids(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
	}
	ship := &Ship{
		PlayerID: 0,
		Position: Position{X: 0, Y: 0},
	}

	HandleShipConquering(m, ship)
}

func TestHandleShipConquering_TooFar(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Size: 10, OwnerID: 1, OwnedSurface: 50, Position: Position{X: 50, Y: 0}},
		},
	}
	ship := &Ship{
		PlayerID: 0,
		Position: Position{X: 0, Y: 0},
	}

	originalSurface := m.Asteroids[0].OwnedSurface

	HandleShipConquering(m, ship)

	if m.Asteroids[0].OwnedSurface != originalSurface {
		t.Errorf("HandleShipConquering() should not conquer when too far, got %v", m.Asteroids[0].OwnedSurface)
	}
}

func TestConquerAsteroid_IncreaseSurface(t *testing.T) {
	m := &Map{}
	ship := &Ship{PlayerID: 0}
	asteroid := &Asteroid{
		Size:         10,
		OwnerID:      0,
		OwnedSurface: 50,
	}

	totalSurface := 10 * 10 * 3.14159265359
	originalSurface := asteroid.OwnedSurface

	ConquerAsteroid(m, ship, asteroid)

	if asteroid.OwnedSurface <= originalSurface {
		t.Errorf("ConquerAsteroid() OwnedSurface should increase, got %v", asteroid.OwnedSurface)
	}
	if asteroid.OwnedSurface > totalSurface {
		t.Errorf("ConquerAsteroid() OwnedSurface should not exceed total surface, got %v", asteroid.OwnedSurface)
	}
}

func TestConquerAsteroid_DecreaseSurface(t *testing.T) {
	m := &Map{}
	ship := &Ship{PlayerID: 1}
	asteroid := &Asteroid{
		Size:         10,
		OwnerID:      0,
		OwnedSurface: 50,
	}

	originalSurface := asteroid.OwnedSurface

	ConquerAsteroid(m, ship, asteroid)

	if asteroid.OwnedSurface >= originalSurface {
		t.Errorf("ConquerAsteroid() OwnedSurface should decrease, got %v", asteroid.OwnedSurface)
	}
}

func TestConquerAsteroid_OwnerChange(t *testing.T) {
	m := &Map{}
	ship := &Ship{PlayerID: 1}
	asteroid := &Asteroid{
		Size:         10,
		OwnerID:      0,
		OwnedSurface: 1,
	}

	ConquerAsteroid(m, ship, asteroid)

	if asteroid.OwnerID != 1 {
		t.Errorf("ConquerAsteroid() OwnerID should change to 1, got %v", asteroid.OwnerID)
	}
}

func TestTickPlayerShips_Movement(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Position: Position{X: 0, Y: 0}, Vector: Position{X: 5, Y: 10}},
		},
		Wormholes: []*Wormhole{},
	}
	p := &Player{ID: 0}

	TickPlayerShips(m, p)

	if m.Ships[0].Position.X != 5 || m.Ships[0].Position.Y != 10 {
		t.Errorf("TickPlayerShips() Position = %v, want {5, 10}", m.Ships[0].Position)
	}
}

func TestTickPlayerShips_MiningAndConquering(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: DrillShip, Position: Position{X: 0, Y: 0}, Vector: Position{X: 0, Y: 0}},
		},
		Asteroids: []*Asteroid{
			{ID: 0, Type: RockAsteroid, Size: 10, Position: Position{X: 5, Y: 0}},
		},
		Wormholes: []*Wormhole{},
	}
	p := &Player{ID: 0}

	TickPlayerShips(m, p)

	if m.Ships[0].Rock == 0 {
		t.Errorf("TickPlayerShips() should mine, got Rock = %v", m.Ships[0].Rock)
	}
}
