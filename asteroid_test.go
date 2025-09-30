package main

import (
	"math"
	"testing"

	"github.com/aquilax/go-perlin"
)

func TestNewAsteroid(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
		Radius:    100,
	}

	asteroid := NewAsteroid(m)

	if asteroid.ID != 0 {
		t.Errorf("NewAsteroid() ID = %v, want 0", asteroid.ID)
	}
	if asteroid.Position.X < -m.Radius || asteroid.Position.X > m.Radius {
		t.Errorf("NewAsteroid() Position X = %v, want between %v and %v", asteroid.Position.X, -m.Radius, m.Radius)
	}
	if asteroid.Position.Y < -m.Radius || asteroid.Position.Y > m.Radius {
		t.Errorf("NewAsteroid() Position Y = %v, want between %v and %v", asteroid.Position.Y, -m.Radius, m.Radius)
	}
	if asteroid.Type != RockAsteroid && asteroid.Type != FuelAsteroid {
		t.Errorf("NewAsteroid() Type = %v, want RockAsteroid or FuelAsteroid", asteroid.Type)
	}
	if asteroid.Size < MinAsteroidSize || asteroid.Size > MaxAsteroidSize {
		t.Errorf("NewAsteroid() Size = %v, want between %v and %v", asteroid.Size, MinAsteroidSize, MaxAsteroidSize)
	}
	if asteroid.OwnerID != -1 {
		t.Errorf("NewAsteroid() OwnerID = %v, want -1", asteroid.OwnerID)
	}
	if asteroid.OwnedSurface != 0 {
		t.Errorf("NewAsteroid() OwnedSurface = %v, want 0", asteroid.OwnedSurface)
	}
	if len(m.Asteroids) != 1 {
		t.Errorf("NewAsteroid() map asteroids count = %v, want 1", len(m.Asteroids))
	}
}

func TestNewAsteroid_MultipleAsteroids(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
		Radius:    100,
	}

	asteroid1 := NewAsteroid(m)
	asteroid2 := NewAsteroid(m)

	if asteroid1.ID != 0 {
		t.Errorf("First asteroid ID = %v, want 0", asteroid1.ID)
	}
	if asteroid2.ID != 1 {
		t.Errorf("Second asteroid ID = %v, want 1", asteroid2.ID)
	}
	if len(m.Asteroids) != 2 {
		t.Errorf("Map asteroids count = %v, want 2", len(m.Asteroids))
	}
}

func TestNewAsteroidFromShip_RockAsteroid(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
		Radius:    100,
	}
	ship := &Ship{
		PlayerID: 0,
		Position: Position{X: 10, Y: 20},
		Rock:     50,
		Fuel:     30,
	}

	asteroid := NewAsteroidFromShip(m, ship, RockAsteroid)

	expectedSize := math.Sqrt(float64(ship.Rock) / MaterialToSurfaceRatio / math.Pi)
	if math.Abs(asteroid.Size-expectedSize) > 1e-10 {
		t.Errorf("NewAsteroidFromShip() Size = %v, want %v", asteroid.Size, expectedSize)
	}
	if asteroid.Type != RockAsteroid {
		t.Errorf("NewAsteroidFromShip() Type = %v, want RockAsteroid", asteroid.Type)
	}
	if asteroid.OwnerID != ship.PlayerID {
		t.Errorf("NewAsteroidFromShip() OwnerID = %v, want %v", asteroid.OwnerID, ship.PlayerID)
	}
	expectedSurface := asteroid.Size * asteroid.Size * math.Pi
	if math.Abs(asteroid.OwnedSurface-expectedSurface) > 1e-10 {
		t.Errorf("NewAsteroidFromShip() OwnedSurface = %v, want %v", asteroid.OwnedSurface, expectedSurface)
	}
	distance := ship.Position.Distance(asteroid.Position)
	if distance > AsteroidSpawnOffset {
		t.Errorf("NewAsteroidFromShip() distance = %v, want <= %v", distance, AsteroidSpawnOffset)
	}
}

func TestNewAsteroidFromShip_FuelAsteroid(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
		Radius:    100,
	}
	ship := &Ship{
		PlayerID: 0,
		Position: Position{X: 10, Y: 20},
		Rock:     50,
		Fuel:     30,
	}

	asteroid := NewAsteroidFromShip(m, ship, FuelAsteroid)

	expectedSize := math.Sqrt(ship.Fuel / MaterialToSurfaceRatio / math.Pi)
	if math.Abs(asteroid.Size-expectedSize) > 1e-10 {
		t.Errorf("NewAsteroidFromShip() Size = %v, want %v", asteroid.Size, expectedSize)
	}
	if asteroid.Type != FuelAsteroid {
		t.Errorf("NewAsteroidFromShip() Type = %v, want FuelAsteroid", asteroid.Type)
	}
}

func TestMineAsteroid_RockAsteroid(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
	}
	asteroid := &Asteroid{
		ID:           0,
		Type:         RockAsteroid,
		Size:         10,
		OwnerID:      -1,
		OwnedSurface: 0,
	}
	m.Asteroids = append(m.Asteroids, asteroid)

	ship := &Ship{
		Rock: 0,
		Fuel: 50,
	}

	initialMaterial := asteroid.Size * asteroid.Size * math.Pi * MaterialToSurfaceRatio
	MineAsteroid(m, ship, asteroid)

	if ship.Rock != ShipMiningAmount {
		t.Errorf("MineAsteroid() ship Rock = %v, want %v", ship.Rock, ShipMiningAmount)
	}
	if ship.Fuel != 50 {
		t.Errorf("MineAsteroid() ship Fuel = %v, want 50", ship.Fuel)
	}

	newMaterial := asteroid.Size * asteroid.Size * math.Pi * MaterialToSurfaceRatio
	expectedNewMaterial := initialMaterial - float64(ShipMiningAmount)
	if math.Abs(newMaterial-expectedNewMaterial) > 1e-10 {
		t.Errorf("MineAsteroid() new material = %v, want %v", newMaterial, expectedNewMaterial)
	}
}

func TestMineAsteroid_FuelAsteroid(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
	}
	asteroid := &Asteroid{
		ID:           0,
		Type:         FuelAsteroid,
		Size:         10,
		OwnerID:      -1,
		OwnedSurface: 0,
	}
	m.Asteroids = append(m.Asteroids, asteroid)

	ship := &Ship{
		Rock: 0,
		Fuel: 50,
	}

	MineAsteroid(m, ship, asteroid)

	if ship.Rock != 0 {
		t.Errorf("MineAsteroid() ship Rock = %v, want 0", ship.Rock)
	}
	if ship.Fuel != 50+float64(ShipMiningAmount) {
		t.Errorf("MineAsteroid() ship Fuel = %v, want %v", ship.Fuel, 50+float64(ShipMiningAmount))
	}
}

func TestMineAsteroid_Depletion(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
	}
	// Create a very small asteroid that will be completely depleted
	smallSize := 0.1
	asteroid := &Asteroid{
		ID:           0,
		Type:         RockAsteroid,
		Size:         smallSize,
		OwnerID:      -1,
		OwnedSurface: 0,
	}
	m.Asteroids = append(m.Asteroids, asteroid)

	ship := &Ship{
		Rock: 0,
		Fuel: 50,
	}

	MineAsteroid(m, ship, asteroid)

	if m.Asteroids[0] != nil {
		t.Errorf("MineAsteroid() asteroid not nil after depletion")
	}
}

func TestMineAsteroid_PartialMining(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
	}
	asteroid := &Asteroid{
		ID:           0,
		Type:         RockAsteroid,
		Size:         2,
		OwnerID:      1,
		OwnedSurface: 10,
	}
	m.Asteroids = append(m.Asteroids, asteroid)

	ship := &Ship{
		Rock: 0,
		Fuel: 50,
	}

	MineAsteroid(m, ship, asteroid)

	if asteroid.OwnedSurface <= 0 {
		t.Errorf("MineAsteroid() OwnedSurface should be positive, got %v", asteroid.OwnedSurface)
	}
}

func TestUpdateAsteroidPositions(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{},
		Round:     0,
		Radius:    100,
	}
	// Initialize perlin noise generator to avoid nil pointer dereference
	m.perlin = perlin.NewPerlin(2, 2, 3, 12345)

	asteroid1 := &Asteroid{
		ID:       0,
		Position: Position{X: 10, Y: 20},
	}
	asteroid2 := &Asteroid{
		ID:       1,
		Position: Position{X: 30, Y: 40},
	}

	m.Asteroids = append(m.Asteroids, asteroid1, asteroid2)

	originalPos1 := asteroid1.Position
	originalPos2 := asteroid2.Position

	UpdateAsteroidPositions(m)

	if asteroid1.Position.X == originalPos1.X && asteroid1.Position.Y == originalPos1.Y {
		t.Errorf("UpdateAsteroidPositions() asteroid1 position did not change")
	}
	if asteroid2.Position.X == originalPos2.X && asteroid2.Position.Y == originalPos2.Y {
		t.Errorf("UpdateAsteroidPositions() asteroid2 position did not change")
	}
}

func TestUpdateAsteroidPositions_WithNilAsteroids(t *testing.T) {
	m := &Map{
		Asteroids: []*Asteroid{
			{ID: 0, Position: Position{X: 10, Y: 20}},
			nil,
			{ID: 2, Position: Position{X: 30, Y: 40}},
		},
		Round:  0,
		Radius: 100,
	}
	// Initialize perlin noise generator to avoid nil pointer dereference
	m.perlin = perlin.NewPerlin(2, 2, 3, 12345)

	UpdateAsteroidPositions(m)

	if m.Asteroids[0].Position.X == 10 && m.Asteroids[0].Position.Y == 20 {
		t.Errorf("UpdateAsteroidPositions() asteroid0 position did not change")
	}
	if m.Asteroids[1] != nil {
		t.Errorf("UpdateAsteroidPositions() nil asteroid became non-nil")
	}
	if m.Asteroids[2].Position.X == 30 && m.Asteroids[2].Position.Y == 40 {
		t.Errorf("UpdateAsteroidPositions() asteroid2 position did not change")
	}
}
