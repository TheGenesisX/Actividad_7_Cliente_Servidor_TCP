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
var InServerFlag chan int

// Imprimir ...
func Imprimir(proc *Proceso, InServerFlag chan int) {
	for {
		select {
		case id := <-InServerFlag:
			if id == proc.ID {
				return
			} else {
				fmt.Printf("id %d: %d \n", proc.ID, proc.Step)
				proc.Step++
			}
		default:
			fmt.Printf("id %d: %d \n", proc.ID, proc.Step)
			proc.Step++
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func servidor() {
	procs = make([]Proceso, 5)
	for i := int(0); i < 5; i++ {
		procs[i] = Proceso{i, 0, true}
		go Imprimir(&procs[i], InServerFlag)
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

	fmt.Println(procs)

	for i := 0; i < len(procs); i++ {
		fmt.Println("IDENTIFICADOR: ", i)
		if procs[i].InServer == true {
			procs[i].InServer = false
			sendToClient = i
			for _ = range procs {
				InServerFlag <- i
			} // We get it out of the server to give it to the client.
			break
		}
	}

	fmt.Println("PROCESO A ENVIAR: ", sendToClient)

	err := gob.NewEncoder(client).Encode(procs[sendToClient])
	if err != nil {
		fmt.Println(err)
		procs[sendToClient].InServer = true
		go Imprimir(&procs[sendToClient], InServerFlag)
		//   procs[sendToClient].InServer = true
	}
	//   go proceso(&procesos[proceso_a_enviar])
	//   return
	// } else {
	// //   estado_de_procesos[proceso_a_enviar] = false
	// }
	return
}

func continueInServer(proc *Proceso) {
	procs[proc.ID].InServer = true
	go Imprimir(&procs[proc.ID], InServerFlag)
	return
}

func main() {
	go servidor()

	var input string
	fmt.Scanln(&input)
}
