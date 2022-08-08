package main

import (
	"math"
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

func (bandit *Bandit) getPullsNumber(alpha float64) int {
	// Вычисление количества дерганий за руку
	// Необходимое для получения хотя бы одной награды
	// С вероятностью [bandit.probability] в [(1-alpha)]% случаев
	var q float64 = float64(1 - bandit.probability)
	var pullsNumber int = int(math.Log(alpha) / math.Log(q))
	if pullsNumber == 0 {
		pullsNumber++
	}
	return pullsNumber
}

func (bandit *Bandit) getPosterior() float32 {
	return float32(bandit.wins) / float32(bandit.games)
}

func (bandit *Bandit) getGamesReward() float32 {
	return bandit.reward * bandit.getPosterior()
}

func (bandit *Bandit) getDiscountedReward() float32 {
	return bandit.reward * bandit.probability
}
