package main

import (
	"fmt"
	"sync"
)

func main() {
	// startslide OMIT
	var once sync.Once // HL
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody) // HL
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	// endslide OMIT
}