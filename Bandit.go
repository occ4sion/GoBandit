package main

import (
	"math/rand"
)

type Bandit struct {
	id, wins, games     int
	probability, reward float32
}

func (bandit *Bandit) pull() bool {
	bandit.games += 1
	var win bool = rand.Float32() < bandit.probability
	if win {
		bandit.wins += 1
	}
	return win
}
