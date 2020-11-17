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
		fmt.Printf("id %d: %d \n", proc.ID, proc.Step)
		proc.Step++
		time.Sleep(time.Millisecond * 500)
	}
}

func cliente() {
	proc := Proceso{0, 0, false}
	var fromServer Proceso

	client, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	err2 := gob.NewEncoder(client).Encode(proc)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	err3 := gob.NewDecoder(client).Decode(&fromServer)
	if err3 != nil {
		fmt.Println(err3)
		return
	}

	go imprimir(&fromServer)
	client.Close()
}

func main() {
	go cliente()

	var input string
	fmt.Scanln(&input)
}
