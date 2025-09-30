package main

import (
	"testing"
)

func TestNewShip(t *testing.T) {
	m := &Map{
		Ships:  []*Ship{},
		Radius: 100,
	}
	p := &Player{
		ID: 0,
		MotherShip: &Ship{
			Position: Position{X: 10, Y: 20},
		},
	}

	ship := NewShip(m, p, BattleShip)

	if ship.ID != 0 {
		t.Errorf("NewShip() ID = %v, want 0", ship.ID)
	}
	if ship.PlayerID != p.ID {
		t.Errorf("NewShip() PlayerID = %v, want %v", ship.PlayerID, p.ID)
	}
	if ship.Position != p.MotherShip.Position {
		t.Errorf("NewShip() Position = %v, want %v", ship.Position, p.MotherShip.Position)
	}
	if ship.Health != ShipMaxHealth {
		t.Errorf("NewShip() Health = %v, want %v", ship.Health, ShipMaxHealth)
	}
	if ship.Fuel != ShipStartFuel {
		t.Errorf("NewShip() Fuel = %v, want %v", ship.Fuel, ShipStartFuel)
	}
	if ship.Type != BattleShip {
		t.Errorf("NewShip() Type = %v, want %v", ship.Type, BattleShip)
	}
	if ship.Rock != 0 {
		t.Errorf("NewShip() Rock = %v, want 0", ship.Rock)
	}
	if len(m.Ships) != 1 {
		t.Errorf("NewShip() map ships count = %v, want 1", len(m.Ships))
	}
	if m.Ships[0] != ship {
		t.Errorf("NewShip() map ship[0] != created ship")
	}
}

func TestNewShip_MultipleShips(t *testing.T) {
	m := &Map{
		Ships:  []*Ship{},
		Radius: 100,
	}
	p := &Player{
		ID: 0,
		MotherShip: &Ship{
			Position: Position{X: 10, Y: 20},
		},
	}

	ship1 := NewShip(m, p, BattleShip)
	ship2 := NewShip(m, p, DrillShip)

	if ship1.ID != 0 {
		t.Errorf("First ship ID = %v, want 0", ship1.ID)
	}
	if ship2.ID != 1 {
		t.Errorf("Second ship ID = %v, want 1", ship2.ID)
	}
	if len(m.Ships) != 2 {
		t.Errorf("Map ships count = %v, want 2", len(m.Ships))
	}
}

func TestNewShip_DifferentShipTypes(t *testing.T) {
	m := &Map{
		Ships:  []*Ship{},
		Radius: 100,
	}
	p := &Player{
		ID: 0,
		MotherShip: &Ship{
			Position: Position{X: 10, Y: 20},
		},
	}

	shipTypes := []ShipType{MotherShip, SuckerShip, DrillShip, TankerShip, TruckShip, BattleShip}

	for _, shipType := range shipTypes {
		ship := NewShip(m, p, shipType)
		if ship.Type != shipType {
			t.Errorf("NewShip() Type = %v, want %v", ship.Type, shipType)
		}
	}
}

func TestShip_Initialization_DefaultValues(t *testing.T) {
	m := &Map{
		Ships:  []*Ship{},
		Radius: 100,
	}
	p := &Player{
		ID: 0,
		MotherShip: &Ship{
			Position: Position{X: 10, Y: 20},
		},
	}

	ship := NewShip(m, p, SuckerShip)

	if ship.Vector.X != 0 || ship.Vector.Y != 0 {
		t.Errorf("NewShip() Vector = %v, want (0, 0)", ship.Vector)
	}
	if ship.Rock != 0 {
		t.Errorf("NewShip() Rock = %v, want 0", ship.Rock)
	}
}

func TestShip_Initialization_WithExistingShips(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0},
			{ID: 1},
			{ID: 2},
		},
		Radius: 100,
	}
	p := &Player{
		ID: 0,
		MotherShip: &Ship{
			Position: Position{X: 10, Y: 20},
		},
	}

	ship := NewShip(m, p, TankerShip)

	if ship.ID != 3 {
		t.Errorf("NewShip() ID = %v, want 3", ship.ID)
	}
	if len(m.Ships) != 4 {
		t.Errorf("Map ships count = %v, want 4", len(m.Ships))
	}
}
