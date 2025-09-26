package main

type Player struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	RockAmount int    `json:"rock"`
	FuelAmount int    `json:"fuel"`
	MotherShip *Ship  `json:"-"`
	Alive      bool   `json:"alive"`
}

func NewPlayer(m *Map, name string) *Player {
	p := &Player{
		ID:         len(m.Players),
		Name:       name,
		Color:      "white",
		RockAmount: PlayerStartRock,
		FuelAmount: PlayerStartFuel,
		Alive:      true,
	}

	s := &Ship{
		ID:       len(m.Ships),
		PlayerID: p.ID,
		Position: RandomPosition(m),
		Type:     MotherShip,
	}
	p.MotherShip = s

	m.Ships = append(m.Ships, s)
	m.Players = append(m.Players, p)
	return p
}
