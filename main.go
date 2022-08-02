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
		fmt.Printf("~~~ Бюджет на каждый раунд (минимум $%d): ", number)
		fmt.Scan(&budget)
	}

	timeit(true)

	bandits = make([]Bandit, number)

	go makeBandits(number)

	timeit()

	number = solve(budget)
	if number == -1 {
		fmt.Println("Бюджет оказался слишком мал")
	} else {
		fmt.Println("Максимальный выигрыш даст бандит №", number)
		getInfo(0.05, number)
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

func solve(budget int) int {
	var bandit *Bandit

	var weights []float32 = make([]float32, len(bandits))
	var winners []bool = make([]bool, len(bandits))

	var maxDiscountSum, maxSum, Sum float32
	count := len(winners)

	for index := range winners {
		winners[index] = true
		Sum += bandits[index].reward
	}
	for index := range weights {
		weights[index] = bandits[index].reward / Sum
	}

	Sum = 0

	for count > 1 {

		pulls := budget / count

		for index := range winners {
			bandit = &bandits[index]

			for i := 0; i < pulls; i++ {
				bandit.pull()
			}

			Sum = weights[index] * bandit.getDiscountedReward()
			if Sum > maxSum {
				maxSum = Sum
			}
			if Sum < maxDiscountSum {
				winners[index] = false
				count -= 1
			}
		}
		maxDiscountSum, maxSum = maxSum, 0
		println()
	}
	var winner int = -1
	for index, win := range winners {
		if win {
			winner = index
		}
	}
	return winner
}
