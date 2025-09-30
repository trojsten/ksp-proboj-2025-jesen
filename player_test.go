package main

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	m := &Map{
		Players: []*Player{},
		Ships:   []*Ship{},
		Radius:  100,
	}
	playerName := "TestPlayer"

	player := NewPlayer(m, playerName)

	if player.ID != 0 {
		t.Errorf("NewPlayer() ID = %v, want 0", player.ID)
	}
	if player.Name != playerName {
		t.Errorf("NewPlayer() Name = %v, want %v", player.Name, playerName)
	}
	if player.Color != "white" {
		t.Errorf("NewPlayer() Color = %v, want white", player.Color)
	}
	if player.RockAmount != PlayerStartRock {
		t.Errorf("NewPlayer() RockAmount = %v, want %v", player.RockAmount, PlayerStartRock)
	}
	if player.FuelAmount != PlayerStartFuel {
		t.Errorf("NewPlayer() FuelAmount = %v, want %v", player.FuelAmount, PlayerStartFuel)
	}
	if !player.Alive {
		t.Errorf("NewPlayer() Alive = %v, want true", player.Alive)
	}
	if player.MotherShip == nil {
		t.Errorf("NewPlayer() MotherShip = nil, want non-nil")
	}
	if len(m.Players) != 1 {
		t.Errorf("NewPlayer() map players count = %v, want 1", len(m.Players))
	}
	if m.Players[0] != player {
		t.Errorf("NewPlayer() map player[0] != created player")
	}
	if len(m.Ships) != 1 {
		t.Errorf("NewPlayer() map ships count = %v, want 1", len(m.Ships))
	}
}

func TestNewPlayer_MultiplePlayers(t *testing.T) {
	m := &Map{
		Players: []*Player{},
		Ships:   []*Ship{},
		Radius:  100,
	}

	player1 := NewPlayer(m, "Player1")
	player2 := NewPlayer(m, "Player2")

	if player1.ID != 0 {
		t.Errorf("First player ID = %v, want 0", player1.ID)
	}
	if player2.ID != 1 {
		t.Errorf("Second player ID = %v, want 1", player2.ID)
	}
	if len(m.Players) != 2 {
		t.Errorf("Map players count = %v, want 2", len(m.Players))
	}
	if len(m.Ships) != 2 {
		t.Errorf("Map ships count = %v, want 2", len(m.Ships))
	}
}

func TestNewPlayer_MotherShipProperties(t *testing.T) {
	m := &Map{
		Players: []*Player{},
		Ships:   []*Ship{},
		Radius:  100,
	}

	player := NewPlayer(m, "TestPlayer")

	if player.MotherShip.PlayerID != player.ID {
		t.Errorf("MotherShip PlayerID = %v, want %v", player.MotherShip.PlayerID, player.ID)
	}
	if player.MotherShip.Type != MotherShip {
		t.Errorf("MotherShip Type = %v, want %v", player.MotherShip.Type, MotherShip)
	}
	if player.MotherShip.Position.X < -m.Radius || player.MotherShip.Position.X > m.Radius {
		t.Errorf("MotherShip Position X = %v, want between %v and %v", player.MotherShip.Position.X, -m.Radius, m.Radius)
	}
	if player.MotherShip.Position.Y < -m.Radius || player.MotherShip.Position.Y > m.Radius {
		t.Errorf("MotherShip Position Y = %v, want between %v and %v", player.MotherShip.Position.Y, -m.Radius, m.Radius)
	}
}

func TestPlayer_Resources(t *testing.T) {
	m := &Map{
		Players: []*Player{},
		Ships:   []*Ship{},
		Radius:  100,
	}

	player := NewPlayer(m, "TestPlayer")

	if player.RockAmount != PlayerStartRock {
		t.Errorf("Initial RockAmount = %v, want %v", player.RockAmount, PlayerStartRock)
	}
	if player.FuelAmount != PlayerStartFuel {
		t.Errorf("Initial FuelAmount = %v, want %v", player.FuelAmount, PlayerStartFuel)
	}
}

func TestNewPlayer_WithExistingPlayers(t *testing.T) {
	m := &Map{
		Players: []*Player{
			{ID: 0},
			{ID: 1},
			{ID: 2},
		},
		Ships: []*Ship{
			{ID: 0},
			{ID: 1},
			{ID: 2},
		},
		Radius: 100,
	}

	player := NewPlayer(m, "NewPlayer")

	if player.ID != 3 {
		t.Errorf("NewPlayer() ID = %v, want 3", player.ID)
	}
	if len(m.Players) != 4 {
		t.Errorf("Map players count = %v, want 4", len(m.Players))
	}
	if len(m.Ships) != 4 {
		t.Errorf("Map ships count = %v, want 4", len(m.Ships))
	}
}

func TestNewPlayer_DifferentNames(t *testing.T) {
	m := &Map{
		Players: []*Player{},
		Ships:   []*Ship{},
		Radius:  100,
	}

	testNames := []string{"Alice", "Bob", "Charlie", "Delta", "Echo"}

	for i, name := range testNames {
		player := NewPlayer(m, name)
		if player.Name != name {
			t.Errorf("Player %d Name = %v, want %v", i, player.Name, name)
		}
		if player.ID != i {
			t.Errorf("Player %d ID = %v, want %v", i, player.ID, i)
		}
	}
}
