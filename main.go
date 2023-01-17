package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/classic-terra/db-benchmark/simnode"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go program(&wg)
	wg.Wait()
}

func program(wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := simnode.GetNode()
	fmt.Println(err)
}
