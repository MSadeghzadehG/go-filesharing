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
	wg.Add(2)
	go func ()  {
		discovery.StartService(3000, "127.0.0.1", nodes, 100)
		defer wg.Done()
	}()
	s := file.GetFile{
		CheckPort:	3001, 
		GetPort:	3002, 
		Ip:			"127.0.0.1", 
		Directory:	"./filesToShare/",
	}
	go func ()  {
		file.StartService(s, nodes)
		defer wg.Done()
	}()
	wg.Wait()
}