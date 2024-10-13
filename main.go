package main

import (
	"encoding/json"
	"fmt"
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
	sum := make(chan int)

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
				sum <- number.A + number.B
			}
		}()
	}

	go func() {
		wg.Wait()
		close(sum)
	}()

	totalsum := 0
	for s := range sum {
		totalsum += s
	}
	fmt.Println(totalsum)
}
