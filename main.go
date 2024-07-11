package main

import (
	"fmt"
	cache "go-cache/internal"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cache := cache.NewCache()
	cache.LoadKeys()

	stopCh := make(chan struct{})

	go cache.Monitor(stopCh)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

	<-sigch //the program will be blocked here till ctrl+c is pressed

	//this will stop all the goroutines when the stopCh is passed
	close(stopCh)

	//flush all the keys from the memory to disk
	cache.Flush()
	fmt.Println("exiting")

}
