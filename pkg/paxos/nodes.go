package paxos

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Value struct {
	Data string
}

type Message struct {
	From       int
	Type       string // Types: "prepare", "promise", "propose", "accepted"
	ProposalID int
	Value      Value
}

type Node struct {
	Id       int
	Port     int
	Incoming chan Message
	Peers    map[int]string // ID to address mapping
}

func (n *Node) SendMessage(msg Message, addr string) {
	retryCount := 0
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			if retryCount < 3 {
				retryCount++
				time.Sleep(time.Duration(retryCount) * time.Second)
				continue
			}
			fmt.Printf("Failed to send message to %s after retries: %v\n", addr, err)
			return
		}
		defer conn.Close()

		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(msg); err != nil {
			fmt.Printf("Failed to encode message: %v\n", err)
			return
		}
		break
	}
}

func (n *Node) Listen() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", n.Port))
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go n.HandleConnection(conn)
	}
}

func (n *Node) HandleConnection(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	var msg Message
	for decoder.Decode(&msg) == nil {
		n.Incoming <- msg
	}
}

func (n *Node) SendHeartbeat() {
	heartbeatMsg := Message{From: n.Id, Type: "heartbeat"}
	for _, addr := range n.Peers {
		go n.SendMessage(heartbeatMsg, addr)
	}
}

func (n *Node) MonitorHeartbeats() {
	heartbeats := make(map[int]time.Time)
	for msg := range n.Incoming {
		if msg.Type == "heartbeat" {
			heartbeats[msg.From] = time.Now()
		}
	}

	// Periodically check for nodes that have not sent a heartbeat
	for id, lastBeat := range heartbeats {
		if time.Since(lastBeat) > 10*time.Second {
			fmt.Printf("Node %d is considered down.\n", id)
			RemoveNode(id)
		}
	}
}
