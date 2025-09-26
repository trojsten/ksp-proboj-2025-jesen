package main

import (
	"fmt"

	"github.com/trojsten/ksp-proboj/client"
)

func StartGame(runner client.Runner) *Map {
	m := NewMap()
	m.runner = &runner
	playerNames, _ := runner.ReadConfig()

	for _, name := range playerNames {
		NewPlayer(m, name)
	}

	runner.Log(fmt.Sprintf("game ready for %d players", len(m.Players)))

	return m
}
