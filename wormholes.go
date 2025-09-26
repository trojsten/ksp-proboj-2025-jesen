package main

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
