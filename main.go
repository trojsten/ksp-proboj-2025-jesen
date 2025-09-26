package main

import (
	"github.com/trojsten/ksp-proboj/client"
)

func main() {
	runner := client.NewRunner()
	m := StartGame(runner)

	for m.ShouldContinue() {
		GameTick(m)
	}
}
