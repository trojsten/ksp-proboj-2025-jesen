package main

import (
	"testing"

	"github.com/trojsten/ksp-proboj/client"
)

func TestFullGameRound(t *testing.T) {
	m := NewMap()

	NewPlayer(m, "Player1")
	NewPlayer(m, "Player2")

	initialRound := m.Round
	initialAsteroidCount := len(m.Asteroids)

	m.Tick()

	if m.Round != initialRound+1 {
		t.Errorf("FullGameRound() Round = %v, want %v", m.Round, initialRound+1)
	}

	if len(m.Asteroids) != initialAsteroidCount {
		t.Errorf("FullGameRound() Asteroids count = %v, want %v", len(m.Asteroids), initialAsteroidCount)
	}

	if len(m.Players) != 2 {
		t.Errorf("FullGameRound() Players count = %v, want 2", len(m.Players))
	}

	if len(m.Ships) != 2 {
		t.Errorf("FullGameRound() Ships count = %v, want 2", len(m.Ships))
	}
}

func TestShipDestructionAndAsteroidCreation(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: BattleShip, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 1, Type: DrillShip, Health: 20, Position: Position{X: 50, Y: 0}, Rock: 30, Fuel: 40},
		},
		Players: []*Player{
			{ID: 0, MotherShip: &Ship{Position: Position{X: 200, Y: 200}}},
			{ID: 1, MotherShip: &Ship{Position: Position{X: 300, Y: 300}}},
		},
		Asteroids: []*Asteroid{},
		Wormholes: []*Wormhole{},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}

	turn := ShootTurnData{SourceID: 0, DestinationID: 1}
	err := turn.Execute(m, p)

	if err != nil {
		t.Errorf("ShipDestructionAndAsteroidCreation() ShootTurn error = %v", err)
	}

	if m.Ships[1] != nil {
		t.Errorf("ShipDestructionAndAsteroidCreation() destroyed ship should be nil")
	}

	if len(m.Asteroids) != 2 {
		t.Errorf("ShipDestructionAndAsteroidCreation() should create 2 asteroids, got %v", len(m.Asteroids))
	}

	if m.Asteroids[0].Type != FuelAsteroid {
		t.Errorf("ShipDestructionAndAsteroidCreation() first asteroid should be FuelAsteroid, got %v", m.Asteroids[0].Type)
	}

	if m.Asteroids[1].Type != RockAsteroid {
		t.Errorf("ShipDestructionAndAsteroidCreation() second asteroid should be RockAsteroid, got %v", m.Asteroids[1].Type)
	}
}

func TestWormholeTeleportationIntegration(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Position: Position{X: 2, Y: 2}, Vector: Position{X: 1, Y: 0}},
		},
	}

	originalPos := m.Ships[0].Position
	CheckShipWormholeTeleportation(m, m.Ships[0])

	if m.Ships[0].Position.X == originalPos.X && m.Ships[0].Position.Y == originalPos.Y {
		t.Errorf("WormholeTeleportationIntegration() ship should have teleported")
	}

	distance := m.Ships[0].Position.Distance(Position{X: 100, Y: 100})
	if distance > WormholeTeleportDistance+1 {
		t.Errorf("WormholeTeleportationIntegration() ship too far from target, distance = %v", distance)
	}
}

func TestResourceManagement(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: MotherShip, Rock: 200, Fuel: 100, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 0, Type: BattleShip, Rock: 0, Fuel: 50, Position: Position{X: 5, Y: 0}},
			{ID: 2, PlayerID: 0, Type: DrillShip, Rock: 0, Fuel: 30, Position: Position{X: 10, Y: 0}},
		},
		Asteroids: []*Asteroid{
			{ID: 0, Type: RockAsteroid, Size: 20, Position: Position{X: 15, Y: 0}},
			{ID: 1, Type: FuelAsteroid, Size: 20, Position: Position{X: 20, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}

	turn1 := LoadTurnData{SourceID: 0, DestinationID: 1, Amount: 50}
	err := turn1.Execute(m, p)
	if err != nil {
		t.Errorf("ResourceManagement() LoadTurn error = %v", err)
	}

	if m.Ships[0].Rock != 150 {
		t.Errorf("ResourceManagement() mothership Rock = %v, want 150", m.Ships[0].Rock)
	}
	if m.Ships[1].Rock != 50 {
		t.Errorf("ResourceManagement() battleship Rock = %v, want 50", m.Ships[1].Rock)
	}

	// Reset used ships for next turn
	m.UsedShips = make(map[int]map[int]bool)

	turn2 := SiphonTurnData{SourceID: 0, DestinationID: 2, Amount: 20}
	err = turn2.Execute(m, p)
	if err != nil {
		t.Errorf("ResourceManagement() SiphonTurn error = %v", err)
	}

	if m.Ships[0].Fuel != 80 {
		t.Errorf("ResourceManagement() mothership Fuel = %v, want 80", m.Ships[0].Fuel)
	}
	if m.Ships[2].Fuel != 50 {
		t.Errorf("ResourceManagement() drillship Fuel = %v, want 50", m.Ships[2].Fuel)
	}

	HandleShipMining(m, m.Ships[2])
	if m.Ships[2].Rock == 0 {
		t.Errorf("ResourceManagement() drillship should have mined rock, got %v", m.Ships[2].Rock)
	}
}

func TestMultipleTurnExecution(t *testing.T) {
	// Create a real runner to avoid nil pointer dereference
	runner := client.NewRunner()
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: MotherShip, Rock: 300, Fuel: 200, Position: Position{X: 0, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
		runner:    &runner,
	}
	p := &Player{ID: 0, RockAmount: 300, FuelAmount: 200, MotherShip: m.Ships[0]}

	turns := []TurnContainer{
		{
			Type: BuyTurn,
			Data: []byte(`{"type":1}`),
		},
		{
			Type: BuyTurn,
			Data: []byte(`{"type":2}`),
		},
	}

	ExecuteTurns(m, p, turns)

	if len(m.Ships) != 3 {
		t.Errorf("MultipleTurnExecution() Ships count = %v, want 3", len(m.Ships))
	}

	if p.RockAmount != 100 {
		t.Errorf("MultipleTurnExecution() player RockAmount = %v, want 100", p.RockAmount)
	}

	if len(m.Ships) > 1 && m.Ships[1].Type != SuckerShip {
		t.Errorf("MultipleTurnExecution() first new ship Type = %v, want SuckerShip", m.Ships[1].Type)
	}

	if len(m.Ships) > 2 && m.Ships[2].Type != DrillShip {
		t.Errorf("MultipleTurnExecution() second new ship Type = %v, want DrillShip", m.Ships[2].Type)
	}
}

func TestGameLoopWithMultiplePlayers(t *testing.T) {
	m := NewMap()

	player1 := NewPlayer(m, "Player1")
	player2 := NewPlayer(m, "Player2")

	turns1 := []TurnContainer{
		{
			Type: BuyTurn,
			Data: []byte(`{"type":5}`),
		},
	}
	turns2 := []TurnContainer{
		{
			Type: BuyTurn,
			Data: []byte(`{"type":2}`),
		},
	}

	ExecuteTurns(m, player1, turns1)
	ExecuteTurns(m, player2, turns2)

	if len(m.Ships) != 4 {
		t.Errorf("GameLoopWithMultiplePlayers() Ships count = %v, want 4", len(m.Ships))
	}

	if player1.RockAmount != PlayerStartRock-100 {
		t.Errorf("GameLoopWithMultiplePlayers() player1 RockAmount = %v, want %v", player1.RockAmount, PlayerStartRock-100)
	}

	if player2.RockAmount != PlayerStartRock-100 {
		t.Errorf("GameLoopWithMultiplePlayers() player2 RockAmount = %v, want %v", player2.RockAmount, PlayerStartRock-100)
	}
}

func TestShipMovementAndWormholeInteraction(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Position: Position{X: 10, Y: 10}, Vector: Position{X: -8, Y: -8}, Fuel: 100},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}

	moveTurn := MoveTurnData{ShipID: 0, Vector: Position{X: -8, Y: -8}}
	err := moveTurn.Execute(m, p)
	if err != nil {
		t.Errorf("ShipMovementAndWormholeInteraction() MoveTurn error = %v", err)
	}

	TickPlayerShips(m, p)

	if m.Ships[0].Position.X == 2 && m.Ships[0].Position.Y == 2 {
		t.Errorf("ShipMovementAndWormholeInteraction() ship should have been teleported")
	}

	// The ship should have been teleported, so it shouldn't be at original position
	if m.Ships[0].Position.X == 2 && m.Ships[0].Position.Y == 2 {
		t.Errorf("ShipMovementAndWormholeInteraction() ship should have teleported")
	}
}

func TestAsteroidMiningAndConqueringCycle(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: DrillShip, Position: Position{X: 0, Y: 0}, Rock: 0},
			{ID: 1, PlayerID: 1, Type: DrillShip, Position: Position{X: 5, Y: 0}, Rock: 0},
		},
		Asteroids: []*Asteroid{
			{ID: 0, Type: RockAsteroid, Size: 20, Position: Position{X: 8, Y: 0}, OwnerID: -1, OwnedSurface: 0},
		},
		Wormholes: []*Wormhole{},
	}
	HandleShipMining(m, m.Ships[0])
	HandleShipConquering(m, m.Ships[0])

	if m.Ships[0].Rock == 0 {
		t.Errorf("AsteroidMiningAndConqueringCycle() player1 ship should have mined rock")
	}

	if m.Asteroids[0].OwnerID != 0 {
		t.Errorf("AsteroidMiningAndConqueringCycle() asteroid should be owned by player1, got %v", m.Asteroids[0].OwnerID)
	}

	HandleShipMining(m, m.Ships[1])
	HandleShipConquering(m, m.Ships[1])

	if m.Ships[1].Rock == 0 {
		t.Errorf("AsteroidMiningAndConqueringCycle() player2 ship should have mined rock")
	}

	// After one conquering cycle, player2 shouldn't own the asteroid yet
	// but surface should be reduced from original
	if m.Asteroids[0].OwnedSurface >= 50 {
		t.Errorf("AsteroidMiningAndConqueringCycle() asteroid surface should be reduced, got %v", m.Asteroids[0].OwnedSurface)
	}
}
