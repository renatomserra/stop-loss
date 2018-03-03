package main

import (
	"./types"
	"fmt"
	"time"
)

func main() {
	state := types.State{}
	fmt.Println(state)

	start := time.Now()

	state.RefreshAll()

	end := time.Now()

	fmt.Println(state)
	fmt.Printf("Time Elapsed: %v", end.Sub(start))
}
