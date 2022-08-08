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
	bandits = make([]Bandit, number)
	go makeBandits(number)

	for budget <= number {
		fmt.Printf("~~~ Бюджет на каждый раунд (минимум $%d): ", number+1)
		fmt.Scan(&budget)
	}

	timeit(true)

	solutions := map[string]Solver{
		"Обычный невзвешенный метод":      {weighted: false, use_posterior: false},
		"Обычный взвешенный метод":        {weighted: true, use_posterior: false},
		"Постериорный невзвешенный метод": {weighted: false, use_posterior: true},
		"Постериорный взвешенный метод":   {weighted: true, use_posterior: true},
	}
	for name, solver := range solutions {
		solver.init()
		index := solver.solve(budget)
		if index == -1 {
			fmt.Println("Бюджет оказался слишком мал")
		} else {
			fmt.Printf("%s посчитал, что максимальный выигрыш даст бандит №%d\n", name, index)
			getInfo(0.05, index)
		}
	}

	timeit()

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
