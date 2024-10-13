package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
)

type Number struct {
	A int
	B int
}

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("number of goroutines is not specified")
	}

	goRoutines, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile("numbers.json")
	if err != nil {
		panic(err)
	}

	numbers := make([]Number, 0)
	err = json.Unmarshal(file, &numbers)
	if err != nil {
		panic(err)
	}

	tasks := make(chan Number)
	sum := make(chan *big.Int)

	go func() {
		for _, number := range numbers {
			tasks <- number
		}
		close(tasks)
	}()

	var wg sync.WaitGroup
	wg.Add(goRoutines)
	for i := 0; i < goRoutines; i++ {
		go func() {
			defer wg.Done()
			for number := range tasks {
				sum <- big.NewInt(int64(number.A + number.B))
			}
		}()
	}

	go func() {
		wg.Wait()
		close(sum)
	}()

	totalsum := big.NewInt(0)
	for s := range sum {
		totalsum.Add(totalsum, s)
	}
	fmt.Println(totalsum)
}
