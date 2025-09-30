package main

import (
	"math"
	"testing"
)

func TestNewWormholes(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{},
		Radius:    100,
	}

	w1, w2 := NewWormholes(m)

	if w1.ID != 0 {
		t.Errorf("First wormhole ID = %v, want 0", w1.ID)
	}
	if w2.ID != 1 {
		t.Errorf("Second wormhole ID = %v, want 1", w2.ID)
	}
	if w1.TargetID != w2.ID {
		t.Errorf("First wormhole TargetID = %v, want %v", w1.TargetID, w2.ID)
	}
	if w2.TargetID != w1.ID {
		t.Errorf("Second wormhole TargetID = %v, want %v", w2.TargetID, w1.ID)
	}
	if w1.Position.X < -m.Radius || w1.Position.X > m.Radius {
		t.Errorf("First wormhole Position X = %v, want between %v and %v", w1.Position.X, -m.Radius, m.Radius)
	}
	if w1.Position.Y < -m.Radius || w1.Position.Y > m.Radius {
		t.Errorf("First wormhole Position Y = %v, want between %v and %v", w1.Position.Y, -m.Radius, m.Radius)
	}
	if w2.Position.X < -m.Radius || w2.Position.X > m.Radius {
		t.Errorf("Second wormhole Position X = %v, want between %v and %v", w2.Position.X, -m.Radius, m.Radius)
	}
	if w2.Position.Y < -m.Radius || w2.Position.Y > m.Radius {
		t.Errorf("Second wormhole Position Y = %v, want between %v and %v", w2.Position.Y, -m.Radius, m.Radius)
	}
	if len(m.Wormholes) != 2 {
		t.Errorf("Map wormholes count = %v, want 2", len(m.Wormholes))
	}
	if m.Wormholes[0] != w1 {
		t.Errorf("Map wormholes[0] != first wormhole")
	}
	if m.Wormholes[1] != w2 {
		t.Errorf("Map wormholes[1] != second wormhole")
	}
}

func TestNewWormholes_MultiplePairs(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{},
		Radius:    100,
	}

	w1a, w1b := NewWormholes(m)
	w2a, w2b := NewWormholes(m)

	if w1a.ID != 0 {
		t.Errorf("First pair first wormhole ID = %v, want 0", w1a.ID)
	}
	if w1b.ID != 1 {
		t.Errorf("First pair second wormhole ID = %v, want 1", w1b.ID)
	}
	if w2a.ID != 2 {
		t.Errorf("Second pair first wormhole ID = %v, want 2", w2a.ID)
	}
	if w2b.ID != 3 {
		t.Errorf("Second pair second wormhole ID = %v, want 3", w2b.ID)
	}
	if len(m.Wormholes) != 4 {
		t.Errorf("Map wormholes count = %v, want 4", len(m.Wormholes))
	}
}

func TestCheckShipWormholeTeleportation_NilShip(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
	}

	CheckShipWormholeTeleportation(m, nil)
}

func TestCheckShipWormholeTeleportation_NoWormholes(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{},
	}
	ship := &Ship{
		Position: Position{X: 0, Y: 0},
		Vector:   Position{X: 1, Y: 1},
	}

	originalPos := ship.Position
	CheckShipWormholeTeleportation(m, ship)

	if ship.Position.X != originalPos.X || ship.Position.Y != originalPos.Y {
		t.Errorf("Ship position changed when no wormholes exist")
	}
}

func TestCheckShipWormholeTeleportation_OutOfRange(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
	}
	ship := &Ship{
		Position: Position{X: 10, Y: 10},
		Vector:   Position{X: 1, Y: 1},
	}

	originalPos := ship.Position
	CheckShipWormholeTeleportation(m, ship)

	if ship.Position.X != originalPos.X || ship.Position.Y != originalPos.Y {
		t.Errorf("Ship position changed when out of wormhole range")
	}
}

func TestCheckShipWormholeTeleportation_WithVector(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
	}
	ship := &Ship{
		Position: Position{X: 2, Y: 2},
		Vector:   Position{X: 3, Y: 4},
	}

	CheckShipWormholeTeleportation(m, ship)

	expectedDistance := math.Sqrt(3*3 + 4*4)
	expectedX := 100 + (3/expectedDistance)*WormholeTeleportDistance
	expectedY := 100 + (4/expectedDistance)*WormholeTeleportDistance

	if math.Abs(ship.Position.X-expectedX) > 1e-10 {
		t.Errorf("Ship Position X = %v, want %v", ship.Position.X, expectedX)
	}
	if math.Abs(ship.Position.Y-expectedY) > 1e-10 {
		t.Errorf("Ship Position Y = %v, want %v", ship.Position.Y, expectedY)
	}
}

func TestCheckShipWormholeTeleportation_WithoutVector(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
	}
	ship := &Ship{
		Position: Position{X: 2, Y: 2},
		Vector:   Position{X: 0, Y: 0},
	}

	CheckShipWormholeTeleportation(m, ship)

	distance := ship.Position.Distance(Position{X: 100, Y: 100})
	if math.Abs(distance-WormholeTeleportDistance) > 1e-10 {
		t.Errorf("Ship distance from target = %v, want %v", distance, WormholeTeleportDistance)
	}
}

func TestCheckShipWormholeTeleportation_MultipleWormholes(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
			{ID: 2, TargetID: 3, Position: Position{X: 50, Y: 50}},
			{ID: 3, TargetID: 2, Position: Position{X: 150, Y: 150}},
		},
	}
	ship := &Ship{
		Position: Position{X: 52, Y: 52},
		Vector:   Position{X: 1, Y: 0},
	}

	CheckShipWormholeTeleportation(m, ship)

	distance := ship.Position.Distance(Position{X: 150, Y: 150})
	if math.Abs(distance-WormholeTeleportDistance) > 1e-10 {
		t.Errorf("Ship should have teleported to closer wormhole, distance = %v", distance)
	}
}

func TestWormhole_TeleportWithVector(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
	}
	ship := &Ship{
		Position: Position{X: 1, Y: 1},
		Vector:   Position{X: 5, Y: 0},
	}

	CheckShipWormholeTeleportation(m, ship)

	expectedX := 100.0 + WormholeTeleportDistance
	expectedY := 100.0

	if math.Abs(ship.Position.X-expectedX) > 1e-10 {
		t.Errorf("Teleport with vector X = %v, want %v", ship.Position.X, expectedX)
	}
	if math.Abs(ship.Position.Y-expectedY) > 1e-10 {
		t.Errorf("Teleport with vector Y = %v, want %v", ship.Position.Y, expectedY)
	}
}

func TestWormhole_TeleportWithoutVector(t *testing.T) {
	m := &Map{
		Wormholes: []*Wormhole{
			{ID: 0, TargetID: 1, Position: Position{X: 0, Y: 0}},
			{ID: 1, TargetID: 0, Position: Position{X: 100, Y: 100}},
		},
	}
	ship := &Ship{
		Position: Position{X: 1, Y: 1},
		Vector:   Position{X: 0, Y: 0},
	}

	CheckShipWormholeTeleportation(m, ship)

	distance := ship.Position.Distance(Position{X: 100, Y: 100})
	if math.Abs(distance-WormholeTeleportDistance) > 1e-10 {
		t.Errorf("Teleport without vector distance = %v, want %v", distance, WormholeTeleportDistance)
	}
}
