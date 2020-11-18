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

var procs []Proceso

func imprimir(proc *Proceso) {
	for {
		switch proc.InServer {
		case true:
			fmt.Printf("id %d: %d \n", proc.ID, proc.Step)
			proc.Step++
			time.Sleep(time.Millisecond * 500)
		case false:
			return // Terminamos goroutine.
		}
	}
}

func servidor() {
	procs = make([]Proceso, 5)
	for i := int(0); i < 5; i++ {
		procs[i] = Proceso{i, 0, true}
		go imprimir(&procs[i])
	}

	server, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Sprintln(err)
		return
	}

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go clientHandler(client)
	}
}

func clientHandler(client net.Conn) {
	var proceso Proceso

	err := gob.NewDecoder(client).Decode(&proceso)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		switch proceso.Step {
		case 0:
			go start(client)
		default:
			go continueInServer(&proceso)
		}
	}
}

func start(client net.Conn) {
	var sendToClient int

	for i := 0; i < len(procs); i++ {
		if procs[i].InServer == true {
			procs[i].InServer = false
			sendToClient = i
			break
		}
	}

	err := gob.NewEncoder(client).Encode(procs[sendToClient])
	if err != nil {
		fmt.Println(err)
		procs[sendToClient].InServer = true
		go imprimir(&procs[sendToClient])
	}
	return
}

func continueInServer(proc *Proceso) {
	procs[proc.ID].InServer = true
	procs[proc.ID].Step = proc.Step
	go imprimir(&procs[proc.ID])
	return
}

func main() {
	go servidor()

	var input string
	fmt.Scanln(&input)
}
