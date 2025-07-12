package main

import (
	"fmt"

	sapp "github.com/GiffE/gokol/sapp"
)

func main() {
	sapp.Run(&sapp.AppDesc{
		Width:   800,
		Height:  600,
		Init:    func() { fmt.Println("initialize") },
		Cleanup: func() { fmt.Println("cleanup") },
		Frame:   func() {},
		Event:   func(e sapp.Event) {},
	})
}
