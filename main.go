package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var bandits []Bandit
var prev time.Time

func init() {
	rand.Seed(1)
	prev = time.Now()
}

func main() {
	var mutex sync.Mutex
	var number int
	fmt.Print("~~~ Количество слот-машин: ")
	fmt.Scan(&number)

	timeit(true)

	bandits = make([]Bandit, number)

	go func() {
		for index := 0; index < number; index++ {
			mutex.Lock()
			bandits[index] = Bandit{
				id: rand.Intn(10000), wins: 0, games: 0,
				probability: rand.Float32(), reward: float32(rand.Intn(10000)),
			}
			mutex.Unlock()
		}
	}()

	timeit() // 0s vs 1.3 ms for 1E4 bandits

	for {
		fmt.Print("~~~ № бандита: ")
		fmt.Scan(&number)
		if number < len(bandits) && number >= 0 {
			getInfo(number)
		} else {
			break
		}
	}
}

func timeit(init ...bool) {
	if !(len(init) != 0 && init[0] == true) {
		fmt.Println(time.Since(prev))
	}
	prev = time.Now()
}

func getInfo(numbers ...int) {
	if len(numbers) == 0 {
		return
	}
	var bandit Bandit
	for _, number := range numbers {
		bandit = bandits[number]
		fmt.Printf("\n================\n||____%4d____||\n|| [7][7][7]  ||  __\n||____________|| (__)\n||            || //\n|| wins|games ||//\n||%4d |%5d ||/\n||            ||\n||    prob    ||\n||    %.2f    ||\n||            ||\n||   reward   ||\n||   $%6.1f  ||\n||============||\n",
			bandit.id, bandit.wins, bandit.games, bandit.probability, bandit.reward)
	}
}

/*

================
||____8081____||
|| [7][7][7]  ||  __
||____________|| (__)
||            || //
|| wins|games ||//
||   0 |    0 ||/
||            ||
||    prob    ||
||    0.94    ||
||            ||
||   reward   ||
||   $1847.0  ||
||============||

*/
