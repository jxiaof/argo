package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	run()
}

func Loop(msg string) {
	for i := 0; i < 20; i++ {
		fmt.Printf(msg+"%d\n", i)
		time.Sleep(time.Second * 1)
	}
}

func run() {
	go Loop("Look, I can count forever1:")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Printf("You pressed ctrl + C. User interrupted infinite loop.")
		os.Exit(0)
	}()
	Loop("Look, I can count forever2:")
}
