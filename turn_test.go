package main

import (
	"encoding/json"
	"testing"
)

func TestParseTurnData_BuyTurn(t *testing.T) {
	data := BuyTurnData{Type: BattleShip}
	jsonData, _ := json.Marshal(data)
	container := TurnContainer{Type: BuyTurn, Data: jsonData}

	turn, err := ParseTurnData(container)
	if err != nil {
		t.Errorf("ParseTurnData() error = %v", err)
	}

	buyTurn, ok := turn.(BuyTurnData)
	if !ok {
		t.Errorf("ParseTurnData() returned wrong type, want BuyTurnData")
	}
	if buyTurn.Type != BattleShip {
		t.Errorf("ParseTurnData() Type = %v, want %v", buyTurn.Type, BattleShip)
	}
}

func TestParseTurnData_MoveTurn(t *testing.T) {
	data := MoveTurnData{ShipID: 1, Vector: Position{X: 5, Y: 10}}
	jsonData, _ := json.Marshal(data)
	container := TurnContainer{Type: MoveTurn, Data: jsonData}

	turn, err := ParseTurnData(container)
	if err != nil {
		t.Errorf("ParseTurnData() error = %v", err)
	}

	moveTurn, ok := turn.(MoveTurnData)
	if !ok {
		t.Errorf("ParseTurnData() returned wrong type, want MoveTurnData")
	}
	if moveTurn.ShipID != 1 {
		t.Errorf("ParseTurnData() ShipID = %v, want 1", moveTurn.ShipID)
	}
	if moveTurn.Vector.X != 5 || moveTurn.Vector.Y != 10 {
		t.Errorf("ParseTurnData() Vector = %v, want {5, 10}", moveTurn.Vector)
	}
}

func TestParseTurnData_InvalidType(t *testing.T) {
	container := TurnContainer{Type: 99, Data: []byte("{}")}

	_, err := ParseTurnData(container)
	if err == nil {
		t.Errorf("ParseTurnData() should return error for invalid type")
	}
}

func TestUseShip_FirstUse(t *testing.T) {
	m := &Map{UsedShips: make(map[int]map[int]bool)}
	p := &Player{ID: 0}

	err := useShip(m, p, 1)
	if err != nil {
		t.Errorf("useShip() error = %v", err)
	}
	if !m.UsedShips[p.ID][1] {
		t.Errorf("useShip() ship not marked as used")
	}
}

func TestUseShip_AlreadyUsed(t *testing.T) {
	m := &Map{UsedShips: map[int]map[int]bool{0: {1: true}}}
	p := &Player{ID: 0}

	err := useShip(m, p, 1)
	if err == nil {
		t.Errorf("useShip() should return error for already used ship")
	}
}

func TestBuyTurn_Execute_Success(t *testing.T) {
	m := &Map{
		Ships:  []*Ship{},
		Radius: 100,
	}
	p := &Player{
		ID:         0,
		RockAmount: 200,
		MotherShip: &Ship{Position: Position{X: 10, Y: 20}},
	}
	turn := BuyTurnData{Type: BattleShip}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("BuyTurn.Execute() error = %v", err)
	}
	if p.RockAmount != 100 {
		t.Errorf("BuyTurn.Execute() RockAmount = %v, want 100", p.RockAmount)
	}
	if len(m.Ships) != 1 {
		t.Errorf("BuyTurn.Execute() Ships count = %v, want 1", len(m.Ships))
	}
}

func TestBuyTurn_Execute_InvalidShipType(t *testing.T) {
	m := &Map{Ships: []*Ship{}}
	p := &Player{RockAmount: 200}
	turn := BuyTurnData{Type: MotherShip}

	err := turn.Execute(m, p)
	if err == nil {
		t.Errorf("BuyTurn.Execute() should return error for invalid ship type")
	}
}

func TestBuyTurn_Execute_InsufficientRocks(t *testing.T) {
	m := &Map{Ships: []*Ship{}}
	p := &Player{RockAmount: 50}
	turn := BuyTurnData{Type: BattleShip}

	err := turn.Execute(m, p)
	if err == nil {
		t.Errorf("BuyTurn.Execute() should return error for insufficient rocks")
	}
}

func TestMoveTurn_Execute_Success(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Fuel: 100, Vector: Position{X: 0, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := MoveTurnData{ShipID: 0, Vector: Position{X: 5, Y: 0}}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("MoveTurn.Execute() error = %v", err)
	}
	if m.Ships[0].Vector.X != 5 || m.Ships[0].Vector.Y != 0 {
		t.Errorf("MoveTurn.Execute() Vector = %v, want {5, 0}", m.Ships[0].Vector)
	}
}

func TestMoveTurn_Execute_VectorScaling(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Fuel: 10000, Vector: Position{X: 0, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := MoveTurnData{ShipID: 0, Vector: Position{X: 20000, Y: 0}}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("MoveTurn.Execute() error = %v", err)
	}
	if m.Ships[0].Vector.X != ShipMovementMaxSize || m.Ships[0].Vector.Y != 0 {
		t.Errorf("MoveTurn.Execute() Vector = %v, want {%v, 0}", m.Ships[0].Vector, ShipMovementMaxSize)
	}
}

func TestMoveTurn_Execute_InsufficientFuel(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Fuel: 1, Vector: Position{X: 0, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := MoveTurnData{ShipID: 0, Vector: Position{X: 100, Y: 0}}

	err := turn.Execute(m, p)
	if err == nil {
		t.Errorf("MoveTurn.Execute() should return error for insufficient fuel")
	}
}

func TestLoadTurn_Execute_Success(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Rock: 50, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 0, Rock: 10, Position: Position{X: 5, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := LoadTurnData{SourceID: 0, DestinationID: 1, Amount: 20}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("LoadTurn.Execute() error = %v", err)
	}
	if m.Ships[0].Rock != 30 {
		t.Errorf("LoadTurn.Execute() source Rock = %v, want 30", m.Ships[0].Rock)
	}
	if m.Ships[1].Rock != 30 {
		t.Errorf("LoadTurn.Execute() destination Rock = %v, want 30", m.Ships[1].Rock)
	}
}

func TestLoadTurn_Execute_TooFar(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Rock: 50, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 0, Rock: 10, Position: Position{X: 100, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := LoadTurnData{SourceID: 0, DestinationID: 1, Amount: 20}

	err := turn.Execute(m, p)
	if err == nil {
		t.Errorf("LoadTurn.Execute() should return error for ships too far apart")
	}
}

func TestSiphonTurn_Execute_Success(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Fuel: 50, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 0, Fuel: 10, Position: Position{X: 5, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := SiphonTurnData{SourceID: 0, DestinationID: 1, Amount: 20}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("SiphonTurn.Execute() error = %v", err)
	}
	if m.Ships[0].Fuel != 30 {
		t.Errorf("SiphonTurn.Execute() source Fuel = %v, want 30", m.Ships[0].Fuel)
	}
	if m.Ships[1].Fuel != 30 {
		t.Errorf("SiphonTurn.Execute() destination Fuel = %v, want 30", m.Ships[1].Fuel)
	}
}

func TestShootTurn_Execute_Success(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: BattleShip, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 1, Type: DrillShip, Health: 50, Position: Position{X: 50, Y: 0}},
		},
		Players: []*Player{
			{ID: 0, MotherShip: &Ship{Position: Position{X: 200, Y: 200}}},
			{ID: 1, MotherShip: &Ship{Position: Position{X: 300, Y: 300}}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := ShootTurnData{SourceID: 0, DestinationID: 1}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("ShootTurn.Execute() error = %v", err)
	}
	if m.Ships[1].Health != 25 {
		t.Errorf("ShootTurn.Execute() destination Health = %v, want 25", m.Ships[1].Health)
	}
}

func TestShootTurn_Execute_ShipDestruction(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: BattleShip, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 1, Type: DrillShip, Health: 20, Position: Position{X: 50, Y: 0}, Rock: 30, Fuel: 40},
		},
		Players: []*Player{
			{ID: 0, MotherShip: &Ship{Position: Position{X: 100, Y: 100}}},
			{ID: 1, MotherShip: &Ship{Position: Position{X: 200, Y: 200}}},
		},
		Asteroids: []*Asteroid{},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := ShootTurnData{SourceID: 0, DestinationID: 1}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("ShootTurn.Execute() error = %v", err)
	}
	if m.Ships[1] != nil {
		t.Errorf("ShootTurn.Execute() destroyed ship should be nil")
	}
	if len(m.Asteroids) != 2 {
		t.Errorf("ShootTurn.Execute() should create 2 asteroids, got %v", len(m.Asteroids))
	}
}

func TestShootTurn_Execute_NonBattleShip(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Type: DrillShip, Position: Position{X: 0, Y: 0}},
			{ID: 1, PlayerID: 1, Type: DrillShip, Health: 50, Position: Position{X: 50, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{ID: 0}
	turn := ShootTurnData{SourceID: 0, DestinationID: 1}

	err := turn.Execute(m, p)
	if err == nil {
		t.Errorf("ShootTurn.Execute() should return error for non-BattleShip")
	}
}

func TestRepairTurn_Execute_Success(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Health: 50, Position: Position{X: 10, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{
		ID:         0,
		MotherShip: &Ship{Position: Position{X: 0, Y: 0}},
	}
	turn := RepairTurnData{ShipID: 0}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("RepairTurn.Execute() error = %v", err)
	}
	if m.Ships[0].Health != 80 {
		t.Errorf("RepairTurn.Execute() Health = %v, want 80", m.Ships[0].Health)
	}
}

func TestRepairTurn_Execute_MaxHealth(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Health: 90, Position: Position{X: 10, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{
		ID:         0,
		MotherShip: &Ship{Position: Position{X: 0, Y: 0}},
	}
	turn := RepairTurnData{ShipID: 0}

	err := turn.Execute(m, p)
	if err != nil {
		t.Errorf("RepairTurn.Execute() error = %v", err)
	}
	if m.Ships[0].Health != ShipMaxHealth {
		t.Errorf("RepairTurn.Execute() Health = %v, want %v", m.Ships[0].Health, ShipMaxHealth)
	}
}

func TestRepairTurn_Execute_TooFar(t *testing.T) {
	m := &Map{
		Ships: []*Ship{
			{ID: 0, PlayerID: 0, Health: 50, Position: Position{X: 100, Y: 0}},
		},
		UsedShips: make(map[int]map[int]bool),
	}
	p := &Player{
		ID:         0,
		MotherShip: &Ship{Position: Position{X: 0, Y: 0}},
	}
	turn := RepairTurnData{ShipID: 0}

	err := turn.Execute(m, p)
	if err == nil {
		t.Errorf("RepairTurn.Execute() should return error for ship too far from mothership")
	}
}
