package main

import "github.com/bhbosman/govalr/internal"

func main() {
	fxApp, _ := internal.CreateFxApp()
	if fxApp.Err() != nil {
		return
	}
	fxApp.Run()
}
