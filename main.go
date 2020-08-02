package main

import (
	"go-filesharing/discovery"
	"go-filesharing/file"
	"sync"
)

func main() {
	nodes := make(map[string]int)
	nodes["0.0.0.0"] = 1
	nodes["localhost"] = 1
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func ()  {
		discovery.StartService(3000, "127.0.0.1", nodes, 100)
		defer wg.Done()
	}()
	wg.Wait()
}