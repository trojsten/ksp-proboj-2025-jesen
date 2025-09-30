package main

import (
	"math"
	"testing"
)

func TestPosition_Add(t *testing.T) {
	p1 := Position{X: 1, Y: 2}
	p2 := Position{X: 3, Y: 4}
	result := p1.Add(p2)

	if result.X != 4 || result.Y != 6 {
		t.Errorf("Add() = (%v, %v), want (%v, %v)", result.X, result.Y, 4, 6)
	}
}

func TestPosition_Sub(t *testing.T) {
	p1 := Position{X: 5, Y: 7}
	p2 := Position{X: 2, Y: 3}
	result := p1.Sub(p2)

	if result.X != 3 || result.Y != 4 {
		t.Errorf("Sub() = (%v, %v), want (%v, %v)", result.X, result.Y, 3, 4)
	}
}

func TestPosition_Distance(t *testing.T) {
	p1 := Position{X: 0, Y: 0}
	p2 := Position{X: 3, Y: 4}
	distance := p1.Distance(p2)

	if distance != 5 {
		t.Errorf("Distance() = %v, want 5", distance)
	}
}

func TestPosition_Distance_SamePoint(t *testing.T) {
	p1 := Position{X: 1, Y: 1}
	p2 := Position{X: 1, Y: 1}
	distance := p1.Distance(p2)

	if distance != 0 {
		t.Errorf("Distance() = %v, want 0", distance)
	}
}

func TestPosition_Size(t *testing.T) {
	p := Position{X: 3, Y: 4}
	size := p.Size()

	if size != 5 {
		t.Errorf("Size() = %v, want 5", size)
	}
}

func TestPosition_Size_Zero(t *testing.T) {
	p := Position{X: 0, Y: 0}
	size := p.Size()

	if size != 0 {
		t.Errorf("Size() = %v, want 0", size)
	}
}

func TestPosition_Scale(t *testing.T) {
	p := Position{X: 2, Y: 3}
	result := p.Scale(2)

	if result.X != 4 || result.Y != 6 {
		t.Errorf("Scale() = (%v, %v), want (%v, %v)", result.X, result.Y, 4, 6)
	}
}

func TestPosition_Scale_Zero(t *testing.T) {
	p := Position{X: 5, Y: 7}
	result := p.Scale(0)

	if result.X != 0 || result.Y != 0 {
		t.Errorf("Scale() = (%v, %v), want (%v, %v)", result.X, result.Y, 0, 0)
	}
}

func TestPosition_Normalize(t *testing.T) {
	p := Position{X: 3, Y: 4}
	result := p.Normalize()

	expectedSize := 1.0
	actualSize := result.Size()
	if math.Abs(actualSize-expectedSize) > 1e-10 {
		t.Errorf("Normalize() size = %v, want %v", actualSize, expectedSize)
	}

	if result.X != 0.6 || result.Y != 0.8 {
		t.Errorf("Normalize() = (%v, %v), want (%v, %v)", result.X, result.Y, 0.6, 0.8)
	}
}

func TestPosition_Normalize_Zero(t *testing.T) {
	p := Position{X: 0, Y: 0}
	result := p.Normalize()

	if result.X != 0 || result.Y != 0 {
		t.Errorf("Normalize() = (%v, %v), want (%v, %v)", result.X, result.Y, 0, 0)
	}
}

func TestRandomPosition(t *testing.T) {
	m := &Map{Radius: 100}

	for i := 0; i < 100; i++ {
		pos := RandomPosition(m)

		if pos.X < -m.Radius || pos.X > m.Radius {
			t.Errorf("RandomPosition() X = %v, want between %v and %v", pos.X, -m.Radius, m.Radius)
		}
		if pos.Y < -m.Radius || pos.Y > m.Radius {
			t.Errorf("RandomPosition() Y = %v, want between %v and %v", pos.Y, -m.Radius, m.Radius)
		}
	}
}

func TestRandomOffsetPosition(t *testing.T) {
	original := Position{X: 10, Y: 20}
	maxOffset := 5.0

	for i := 0; i < 100; i++ {
		result := RandomOffsetPosition(original, maxOffset)
		distance := original.Distance(result)

		if distance > maxOffset {
			t.Errorf("RandomOffsetPosition() distance = %v, want <= %v", distance, maxOffset)
		}
	}
}

func TestRandomOffsetPosition_ZeroOffset(t *testing.T) {
	original := Position{X: 10, Y: 20}
	result := RandomOffsetPosition(original, 0)

	if result.X != original.X || result.Y != original.Y {
		t.Errorf("RandomOffsetPosition() = (%v, %v), want (%v, %v)", result.X, result.Y, original.X, original.Y)
	}
}
