package main

import (
	"log"
	"net"
)

type Client struct {
	Conn net.Conn
}

type MessageType uint
const (
	ClientConnected MessageType = iota
	NewMessage
	ClientDisconnected
)

type Message struct {
	Type MessageType
	Sender Client
	Message string
}

const (
	PORT = "6969"
)

func server(incoming chan Message){
	clients := make(map[string]Client, 100)
	for {
		msg :=  <- incoming

		client := Client {
			Conn: msg.Sender.Conn,
		}

		switch msg.Type {
			case ClientConnected: {

				log.Printf("Client %s connected", client.Conn.RemoteAddr())
				clients[client.Conn.RemoteAddr().String()] = client
				break
			}

			case NewMessage: {
				for _, v := range clients {
					if v.Conn.RemoteAddr() != client.Conn.RemoteAddr() {
						_, err := v.Conn.Write([]byte(msg.Message))
						if err != nil {
							log.Printf("Error sending to the client %s : %s", msg.Sender.Conn.RemoteAddr(), err)
						}
					}
				}
				break
			}

			case ClientDisconnected: {
				delete(clients, msg.Sender.Conn.RemoteAddr().String())
			}
		}
	}
}

func client(conn net.Conn, incoming chan Message){
	defer conn.Close()

	incoming <- Message{
		Type: ClientConnected,
		Message: "",
		Sender: Client{
			Conn: conn,
		},
	}

	buffer := make([]byte, 512)
	for {

		_, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection! : %s", err)
			return
		}

		msg_str := string(buffer[:])

		log.Printf("%s", msg_str)

		incoming <- Message{
			Type: NewMessage,
			Message: msg_str,
			Sender: Client{
				Conn: conn,
			},
		}
	}
}

func main(){
	ln, err := net.Listen("tcp", ":" + PORT)
	if err != nil {
		// handle error
	}
	log.Printf("Listening port %s ...", PORT)
	incoming := make(chan Message)
	go server(incoming)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go client(conn, incoming)
	}
}

