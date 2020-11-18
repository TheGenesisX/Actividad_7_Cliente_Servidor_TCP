package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

// Proceso ...
type Proceso struct {
	ID       int
	Step     int
	InServer bool
}

func imprimir(proc *Proceso) {
	for {
		switch proc.InServer {
		case false:
			fmt.Printf("id %d: %d \n", proc.ID, proc.Step)
			proc.Step++
			time.Sleep(time.Millisecond * 500)
		case true:
			return
		}
	}
}

func main() {
	proc := Proceso{0, 0, false}
	// var fromServer Proceso

	client, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = gob.NewEncoder(client).Encode(proc)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = gob.NewDecoder(client).Decode(&proc)
	if err != nil {
		fmt.Println(err)
		return
	}

	go imprimir(&proc)
	client.Close()

	var input string
	fmt.Scanln(&input) // Cuando queremos liberar al cliente del proceso que tiene.

	client, err = net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = gob.NewEncoder(client).Encode(proc)
	if err != nil {
		fmt.Println(err)
		return
	}

	client.Close()
}
