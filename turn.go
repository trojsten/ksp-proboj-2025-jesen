package main

import (
	"encoding/json"
	"fmt"
)

type TurnType int

const (
	BuyTurn TurnType = iota
	MoveTurn
	LoadTurn
	SiphonTurn
	ShootTurn
)

type TurnContainer struct {
	Type TurnType        `json:"type"`
	Data json.RawMessage `json:"data"`
}

func ParseTurnData(container TurnContainer) (Turn, error) {
	switch container.Type {
	case BuyTurn:
		var turn BuyTurnData
		err := json.Unmarshal(container.Data, &turn)
		return turn, err
	}

	return nil, fmt.Errorf("unknown turn type: %v", container.Type)
}

func ExecuteTurns(m *Map, p *Player, turns []TurnContainer) {
	for _, container := range turns {
		turn, err := ParseTurnData(container)
		if err != nil {
			m.runner.Log(fmt.Sprintf("could not parse turn '%v': %v", container, err))
			continue
		}

		err = turn.Validate()
		if err != nil {
			m.runner.Log(fmt.Sprintf("invalid turn data '%v': %v", container, err))
			continue
		}

		err = turn.Execute(m, p)
		if err != nil {
			m.runner.Log(fmt.Sprintf("error while executing turn '%v': %v", container, err))
		}
	}
}

type Turn interface {
	Validate() error
	Execute(*Map, *Player) error
}

type BuyTurnData struct {
	Type ShipType
}

func (t BuyTurnData) Validate() error {
	if t.Type <= MotherShip || t.Type > BattleShip {
		return fmt.Errorf("invalid ship type: %v", t.Type)
	}

	return nil
}

func (t BuyTurnData) Execute(m *Map, p *Player) error {
	price := ShipRockPrice(t.Type)
	if p.RockAmount < price {
		return fmt.Errorf("not enough rocks in mothership")
	}

	p.RockAmount -= price
	NewShip(m, p, t.Type)
	return nil
}
