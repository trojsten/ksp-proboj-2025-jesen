package main

import (
	"math"
	"math/rand"
)

type Wormhole struct {
	ID       int      `json:"id"`
	TargetID int      `json:"target_id"`
	Position Position `json:"position"`
}

func NewWormholes(m *Map) (*Wormhole, *Wormhole) {
	w1 := &Wormhole{
		ID:       len(m.Wormholes),
		Position: RandomPosition(m),
	}

	m.Wormholes = append(m.Wormholes, w1)

	w2 := &Wormhole{
		ID:       len(m.Wormholes),
		TargetID: w1.ID,
		Position: RandomPosition(m),
	}

	m.Wormholes = append(m.Wormholes, w2)
	w1.TargetID = w2.ID

	return w1, w2
}

func CheckWormholeTeleportation(m *Map) {
	for _, ship := range m.Ships {
		if ship == nil {
			continue
		}

		for _, wormhole := range m.Wormholes {
			distance := ship.Position.Distance(wormhole.Position)
			if distance < WormholeRadius {
				// Find the target wormhole
				var targetWormhole *Wormhole
				for _, wh := range m.Wormholes {
					if wh.ID == wormhole.TargetID {
						targetWormhole = wh
						break
					}
				}

				if targetWormhole != nil {
					// Calculate teleport position: target wormhole position + scaled vector
					// Scale the vector to ensure minimum distance of WormholeTeleportDistance
					if ship.Vector.Size() > 0 {
						normalizedVector := ship.Vector.Normalize()
						teleportVector := normalizedVector.Scale(WormholeTeleportDistance)
						ship.Position = targetWormhole.Position.Add(teleportVector)
					} else {
						// If ship has no vector, teleport to a random position at minimum distance
						angle := rand.Float64() * 2 * math.Pi
						teleportX := targetWormhole.Position.X + WormholeTeleportDistance*math.Cos(angle)
						teleportY := targetWormhole.Position.Y + WormholeTeleportDistance*math.Sin(angle)
						ship.Position = Position{teleportX, teleportY}
					}
					break // Only teleport once per turn
				}
			}
		}
	}
}
