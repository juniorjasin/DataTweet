package main

import (
	"runtime"
	"fmt"
)

func main() {
	configureRuntime()
  mapURLsToControllers()
}

func configureRuntime(){
	numCPU := runtime.NumCPU()
	fmt.Println(numCPU)
	runtime.GOMAXPROCS(numCPU)
}
