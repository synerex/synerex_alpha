package sxutil

import (
	"log"
	"os"
	"os/signal"
)

// Signal handling Utilities for Synergic Market

var (
	funcSlice []func()
)

func init() {
	funcSlice = make([]func(), 0)
}

// register closing functions.
func RegisterDeferFunction(f func()) {
	funcSlice = append(funcSlice, f)
}

func CallDeferFunctions() {
	for _, f := range funcSlice {
		log.Printf("Calling %v", f)
		f()
	}
}

func HandleSigInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	for range c {
		log.Println("Signal Interrupt!")
		close(c)
	}

	CallDeferFunctions()

	log.Println("End at HandleSigInt in sxutil/signal.go")
	os.Exit(1)
}
