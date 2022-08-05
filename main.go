package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var bandits []Bandit
var prev time.Time
var mutex sync.Mutex

func init() {
	rand.Seed(1)
	prev = time.Now()
}

func main() {

	var number, budget int

	fmt.Print("~~~ Количество бандитов: ")
	fmt.Scan(&number)
	for budget <= number {
		fmt.Printf("~~~ Бюджет на каждый раунд (минимум $%d): ", number+1)
		fmt.Scan(&budget)
	}

	timeit(true)

	bandits = make([]Bandit, number)

	go makeBandits(number)

	timeit()

	solutions := map[string]int{
		"Невзвешенный метод": solve(budget, false),
		"Взвешенный метод":   solve(budget, true),
	}
	for method, index := range solutions {
		if index == -1 {
			fmt.Println("Бюджет оказался слишком мал")
		} else {
			fmt.Printf("%s посчитал, что максимальный выигрыш даст бандит №%d\n", method, index)
			getInfo(0.05, index)
		}
	}

	for {
		fmt.Print("~~~ № бандита: ")
		fmt.Scan(&number)
		if number < len(bandits) && number >= 0 {
			getInfo(0.01, number)
		} else {
			break
		}
	}
}

// ------------------------------------------------------------------------------------------------------

func timeit(init ...bool) {
	if !(len(init) != 0 && init[0] == true) {
		fmt.Println(time.Since(prev))
	}
	prev = time.Now()
}

func getInfo(alpha float64, numbers ...int) {
	if len(numbers) == 0 {
		return
	}
	var bandit Bandit
	for _, number := range numbers {
		bandit = bandits[number]
		fmt.Printf("Дисконтированная выгода = %.2f.\n", bandit.getDiscountedReward())
		fmt.Printf("Понадобится дёрнуть %d раз, чтобы с %.2f%% уверенностью получить выигрыш.\n", bandit.getPullsNumber(alpha), 1-alpha)
		fmt.Printf("================\n||____%4d____||\n|| [7][7][7]  ||  __\n||____________|| (__)\n||            || //\n|| wins|games ||//\n||%4d |%5d ||/\n||            ||\n||    prob    ||\n||    %.2f    ||\n||            ||\n||   reward   ||\n||   $%6.1f  ||\n||============||\n",
			bandit.id, bandit.wins, bandit.games, bandit.probability, bandit.reward)
	}
}

func makeBandits(number int) {
	for index := 0; index < number; index++ {
		mutex.Lock()
		bandits[index] = Bandit{
			id: rand.Intn(10000), wins: 0, games: 0,
			probability: rand.Float32(), reward: float32(rand.Intn(10000)),
		}
		mutex.Unlock()
	}
}

func solve(budget int, weighted bool) int {
	var bandit *Bandit

	var weights []float32 = make([]float32, len(bandits))

	var maxDiscountSum, maxSum, Sum float32

	var winner, round int

	count := len(weights)

	for index := range weights {
		Sum += bandits[index].reward
	}
	if weighted {
		for index := range weights {
			weights[index] = bandits[index].reward / Sum
		}
	} else {
		for index := range weights {
			weights[index] = 1
		}
	}

	Sum = 0
	winner = -1

	for count > 1 {

		pulls := budget / count

		round++

		fmt.Printf("Раунд %d. Бандитов осталось %d.\n", round, count)

		for index := range weights {
			if weights[index] == 0 {
				continue
			}
			Sum = 0

			bandit = &bandits[index]

			for i := 0; i < pulls; i++ {
				if bandit.pull() {
					Sum += bandit.reward
				}
			}

			Sum *= weights[index]
			if Sum > maxSum {
				maxSum = Sum
				winner = index
			}
			if Sum < maxDiscountSum {
				weights[index] = 0
				count -= 1
			}
		}
		maxDiscountSum, maxSum = maxSum, 0
	}
	return winner
}
