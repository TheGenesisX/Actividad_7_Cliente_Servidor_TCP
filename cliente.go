package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

var stop chan bool

// Proceso ...
type Proceso struct {
	ID       int
	Step     int
	InServer bool
}

func Imprimir(proc *Proceso, stop chan bool) {
	for {
		select {
		case <-stop:
			return
		default:
			fmt.Printf("id %d: %d \n", proc.ID, proc.Step)
			proc.Step++
			time.Sleep(time.Millisecond * 500)
		}
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
	// fmt.Println(fromServer)

	go Imprimir(&fromServer, stop)
}

func main() {
	go cliente()

	var input string
	fmt.Scanln(&input)

	// stop <- true
}
