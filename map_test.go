package main

import (
	"testing"

	"github.com/aquilax/go-perlin"
)

func TestNewMap(t *testing.T) {
	m := NewMap()

	if m.Radius != Radius {
		t.Errorf("NewMap() Radius = %v, want %v", m.Radius, Radius)
	}
	if len(m.Ships) != 0 {
		t.Errorf("NewMap() Ships count = %v, want 0", len(m.Ships))
	}
	if len(m.Asteroids) != AsteroidCount {
		t.Errorf("NewMap() Asteroids count = %v, want %v", len(m.Asteroids), AsteroidCount)
	}
	if len(m.Wormholes) != WormholeCount*2 {
		t.Errorf("NewMap() Wormholes count = %v, want %v", len(m.Wormholes), WormholeCount*2)
	}
	if len(m.Players) != 0 {
		t.Errorf("NewMap() Players count = %v, want 0", len(m.Players))
	}
	if m.Round != 0 {
		t.Errorf("NewMap() Round = %v, want 0", m.Round)
	}
	if m.perlin == nil {
		t.Errorf("NewMap() perlin = nil, want non-nil")
	}
	if m.UsedShips != nil {
		t.Errorf("NewMap() UsedShips = %v, want nil", m.UsedShips)
	}
}

func TestNewMap_AsteroidProperties(t *testing.T) {
	m := NewMap()

	for i, asteroid := range m.Asteroids {
		if asteroid.ID != i {
			t.Errorf("Asteroid %d ID = %v, want %d", i, asteroid.ID, i)
		}
		if asteroid.Position.X < -m.Radius || asteroid.Position.X > m.Radius {
			t.Errorf("Asteroid %d Position X = %v, want between %v and %v", i, asteroid.Position.X, -m.Radius, m.Radius)
		}
		if asteroid.Position.Y < -m.Radius || asteroid.Position.Y > m.Radius {
			t.Errorf("Asteroid %d Position Y = %v, want between %v and %v", i, asteroid.Position.Y, -m.Radius, m.Radius)
		}
		if asteroid.Type != RockAsteroid && asteroid.Type != FuelAsteroid {
			t.Errorf("Asteroid %d Type = %v, want RockAsteroid or FuelAsteroid", i, asteroid.Type)
		}
		if asteroid.Size < MinAsteroidSize || asteroid.Size > MaxAsteroidSize {
			t.Errorf("Asteroid %d Size = %v, want between %v and %v", i, asteroid.Size, MinAsteroidSize, MaxAsteroidSize)
		}
		if asteroid.OwnerID != -1 {
			t.Errorf("Asteroid %d OwnerID = %v, want -1", i, asteroid.OwnerID)
		}
		if asteroid.OwnedSurface != 0 {
			t.Errorf("Asteroid %d OwnedSurface = %v, want 0", i, asteroid.OwnedSurface)
		}
	}
}

func TestNewMap_WormholeProperties(t *testing.T) {
	m := NewMap()

	if len(m.Wormholes)%2 != 0 {
		t.Errorf("Wormholes count = %v, want even number", len(m.Wormholes))
	}

	for i := 0; i < len(m.Wormholes); i += 2 {
		w1 := m.Wormholes[i]
		w2 := m.Wormholes[i+1]

		if w1.ID != i {
			t.Errorf("Wormhole %d ID = %v, want %d", i, w1.ID, i)
		}
		if w2.ID != i+1 {
			t.Errorf("Wormhole %d ID = %v, want %d", i+1, w2.ID, i+1)
		}
		if w1.TargetID != w2.ID {
			t.Errorf("Wormhole %d TargetID = %v, want %d", i, w1.TargetID, w2.ID)
		}
		if w2.TargetID != w1.ID {
			t.Errorf("Wormhole %d TargetID = %v, want %d", i+1, w2.TargetID, w1.ID)
		}
		if w1.Position.X < -m.Radius || w1.Position.X > m.Radius {
			t.Errorf("Wormhole %d Position X = %v, want between %v and %v", i, w1.Position.X, -m.Radius, m.Radius)
		}
		if w1.Position.Y < -m.Radius || w1.Position.Y > m.Radius {
			t.Errorf("Wormhole %d Position Y = %v, want between %v and %v", i, w1.Position.Y, -m.Radius, m.Radius)
		}
		if w2.Position.X < -m.Radius || w2.Position.X > m.Radius {
			t.Errorf("Wormhole %d Position X = %v, want between %v and %v", i+1, w2.Position.X, -m.Radius, m.Radius)
		}
		if w2.Position.Y < -m.Radius || w2.Position.Y > m.Radius {
			t.Errorf("Wormhole %d Position Y = %v, want between %v and %v", i+1, w2.Position.Y, -m.Radius, m.Radius)
		}
	}
}

func TestMap_ShouldContinue(t *testing.T) {
	m := &Map{Round: 0}

	if !m.ShouldContinue() {
		t.Errorf("ShouldContinue() = false, want true for round 0")
	}

	m.Round = 500
	if !m.ShouldContinue() {
		t.Errorf("ShouldContinue() = false, want true for round 500")
	}

	m.Round = 1000
	if !m.ShouldContinue() {
		t.Errorf("ShouldContinue() = false, want true for round 1000")
	}

	m.Round = 1001
	if m.ShouldContinue() {
		t.Errorf("ShouldContinue() = true, want false for round 1001")
	}
}

func TestMap_Tick(t *testing.T) {
	m := &Map{
		Round: 5,
		Asteroids: []*Asteroid{
			{ID: 0, Position: Position{X: 10, Y: 20}},
			{ID: 1, Position: Position{X: 30, Y: 40}},
		},
		Radius: 100,
	}
	// Initialize perlin noise generator to avoid nil pointer dereference
	m.perlin = perlin.NewPerlin(2, 2, 3, 12345)

	originalRound := m.Round
	originalPos0 := m.Asteroids[0].Position
	originalPos1 := m.Asteroids[1].Position

	m.Tick()

	if m.Round != originalRound+1 {
		t.Errorf("Tick() Round = %v, want %v", m.Round, originalRound+1)
	}

	if m.Asteroids[0].Position.X == originalPos0.X && m.Asteroids[0].Position.Y == originalPos0.Y {
		t.Errorf("Tick() asteroid0 position did not change")
	}

	if m.Asteroids[1].Position.X == originalPos1.X && m.Asteroids[1].Position.Y == originalPos1.Y {
		t.Errorf("Tick() asteroid1 position did not change")
	}
}

func TestMap_Tick_WithNilAsteroids(t *testing.T) {
	m := &Map{
		Round: 5,
		Asteroids: []*Asteroid{
			{ID: 0, Position: Position{X: 10, Y: 20}},
			nil,
			{ID: 2, Position: Position{X: 30, Y: 40}},
		},
		Radius: 100,
	}
	// Initialize perlin noise generator to avoid nil pointer dereference
	m.perlin = perlin.NewPerlin(2, 2, 3, 12345)

	m.Tick()

	if m.Round != 6 {
		t.Errorf("Tick() Round = %v, want 6", m.Round)
	}

	if m.Asteroids[0].Position.X == 10 && m.Asteroids[0].Position.Y == 20 {
		t.Errorf("Tick() asteroid0 position did not change")
	}

	if m.Asteroids[1] != nil {
		t.Errorf("Tick() nil asteroid became non-nil")
	}

	if m.Asteroids[2].Position.X == 30 && m.Asteroids[2].Position.Y == 40 {
		t.Errorf("Tick() asteroid2 position did not change")
	}
}

func TestMap_Tick_EmptyAsteroids(t *testing.T) {
	m := &Map{
		Round:     5,
		Asteroids: []*Asteroid{},
		Radius:    100,
	}
	// Initialize perlin noise generator to avoid nil pointer dereference
	m.perlin = perlin.NewPerlin(2, 2, 3, 12345)

	m.Tick()

	if m.Round != 6 {
		t.Errorf("Tick() Round = %v, want 6", m.Round)
	}
}
