package main

import "fmt"

type Solver struct {
	use_posterior, weighted bool
	weights                 []float32
}

func (solver *Solver) init() {
	solver.weights = make([]float32, len(bandits))
	if solver.weighted {
		var Sum float32
		for index := range solver.weights {
			Sum += bandits[index].reward
		}
		for index := range solver.weights {
			solver.weights[index] = bandits[index].reward / Sum
		}
	} else {
		for index := range solver.weights {
			solver.weights[index] = 1
		}
	}

}

func (solver *Solver) solve(budget int) int {
	var bandit Bandit

	var maxDiscountSum, maxSum, Sum float32

	var winner, round int

	var count int = len(solver.weights)

	Sum = 0
	winner = -1

	for count > 1 {

		pulls := budget / count

		round++

		fmt.Printf("Раунд %d. Бандитов осталось %d.\n", round, count)

		for index, weight := range solver.weights {
			if weight == 0 {
				continue
			}
			Sum = 0

			bandit = bandits[index]

			for i := 0; i < pulls; i++ {
				if bandit.pull() {
					Sum += bandit.reward
				}
			}

			if solver.use_posterior {
				Sum += Sum*bandit.getPosterior() - Sum
			}
			Sum += Sum*weight - Sum
			if Sum > maxSum {
				maxSum = Sum
				winner = index
			}
			if Sum < maxDiscountSum {
				solver.weights[index] = 0
				count -= 1
			}
		}
		maxDiscountSum, maxSum = maxSum, 0
	}
	return winner
}
