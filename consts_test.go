package main

import (
	"testing"
)

func TestShipRockPrice(t *testing.T) {
	shipTypes := []ShipType{MotherShip, SuckerShip, DrillShip, TankerShip, TruckShip, BattleShip}

	for _, shipType := range shipTypes {
		price := ShipRockPrice(shipType)
		if price != 100 {
			t.Errorf("ShipRockPrice(%v) = %v, want 100", shipType, price)
		}
	}
}

func TestShipMovementPrice_ZeroVector(t *testing.T) {
	vector := Position{X: 0, Y: 0}
	shipType := BattleShip

	price := ShipMovementPrice(vector, shipType)
	if price != 0 {
		t.Errorf("ShipMovementPrice() = %v, want 0 for zero vector", price)
	}
}

func TestShipMovementPrice_FreeMovement(t *testing.T) {
	vector := Position{X: ShipMovementFree, Y: 0}
	shipType := BattleShip

	price := ShipMovementPrice(vector, shipType)
	if price != 0 {
		t.Errorf("ShipMovementPrice() = %v, want 0 for free movement", price)
	}
}

func TestShipMovementPrice_PaidMovement(t *testing.T) {
	vector := Position{X: ShipMovementFree + 5, Y: 0}
	shipType := BattleShip

	expectedPrice := 5.0 * ShipMovementMultiplier
	price := ShipMovementPrice(vector, shipType)

	if price != expectedPrice {
		t.Errorf("ShipMovementPrice() = %v, want %v", price, expectedPrice)
	}
}

func TestShipMovementPrice_DifferentShipTypes(t *testing.T) {
	vector := Position{X: ShipMovementFree + 10, Y: 0}
	shipTypes := []ShipType{MotherShip, SuckerShip, DrillShip, TankerShip, TruckShip, BattleShip}

	for _, shipType := range shipTypes {
		price := ShipMovementPrice(vector, shipType)
		expectedPrice := 10.0 * ShipMovementMultiplier
		if price != expectedPrice {
			t.Errorf("ShipMovementPrice(%v) = %v, want %v", shipType, price, expectedPrice)
		}
	}
}

func TestShipMovementPrice_NegativeVector(t *testing.T) {
	vector := Position{X: -5, Y: -10}
	shipType := BattleShip

	vectorSize := vector.Size()
	expectedPrice := max(0.0, (vectorSize-ShipMovementFree)*ShipMovementMultiplier)
	price := ShipMovementPrice(vector, shipType)
	if price != expectedPrice {
		t.Errorf("ShipMovementPrice() = %v, want %v for negative vector", price, expectedPrice)
	}
}

func TestShipMovementPrice_ComplexVector(t *testing.T) {
	vector := Position{X: 3, Y: 4}
	shipType := BattleShip

	vectorSize := vector.Size()
	if vectorSize < ShipMovementFree {
		t.Errorf("Test setup error: vector size %v should be >= free movement %v", vectorSize, ShipMovementFree)
	}

	expectedPrice := (vectorSize - ShipMovementFree) * ShipMovementMultiplier
	price := ShipMovementPrice(vector, shipType)

	if price != expectedPrice {
		t.Errorf("ShipMovementPrice() = %v, want %v", price, expectedPrice)
	}
}

func TestRandomFloat_BasicRange(t *testing.T) {
	min := 5.0
	max := 10.0

	for i := 0; i < 100; i++ {
		result := RandomFloat(min, max)
		if result < min || result > max {
			t.Errorf("RandomFloat() = %v, want between %v and %v", result, min, max)
		}
	}
}

func TestRandomFloat_SameMinMax(t *testing.T) {
	min := 5.0
	max := 5.0

	result := RandomFloat(min, max)
	if result != min {
		t.Errorf("RandomFloat() = %v, want %v for same min/max", result, min)
	}
}

func TestRandomFloat_NegativeRange(t *testing.T) {
	min := -10.0
	max := -5.0

	for i := 0; i < 100; i++ {
		result := RandomFloat(min, max)
		if result < min || result > max {
			t.Errorf("RandomFloat() = %v, want between %v and %v", result, min, max)
		}
	}
}

func TestRandomFloat_ZeroRange(t *testing.T) {
	min := 0.0
	max := 1.0

	for i := 0; i < 100; i++ {
		result := RandomFloat(min, max)
		if result < min || result > max {
			t.Errorf("RandomFloat() = %v, want between %v and %v", result, min, max)
		}
	}
}

func TestRandomFloat_LargeRange(t *testing.T) {
	min := -1000.0
	max := 1000.0

	for i := 0; i < 100; i++ {
		result := RandomFloat(min, max)
		if result < min || result > max {
			t.Errorf("RandomFloat() = %v, want between %v and %v", result, min, max)
		}
	}
}
