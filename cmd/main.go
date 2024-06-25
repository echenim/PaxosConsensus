package main

import (
	"github.com/echenim/PaxosConsensus/pkg/paxos"
)

func main() {
	// Configuration and peer management
	paxos.AddNode(1, "localhost:8000")
	paxos.AddNode(2, "localhost:8001")
	paxos.AddNode(3, "localhost:8002")
	paxos.AddNode(4, "localhost:8003") // Learner node

	node := &paxos.Node{
		Id:       1,
		Port:     8000,
		Incoming: make(chan paxos.Message, 10),
	}
	go node.Listen()
	// go node.Run()
	go node.SendHeartbeat()
	go node.MonitorHeartbeats()

	// Majority calculation (simple majority)
	majority := (len(paxos.ConfigService.Peers) / 2) + 1

	learner := paxos.NewLearner(4, 8003, paxos.ConfigService.Peers, majority)
	go learner.Listen()
	go learner.Run()

	// other node initializations and operations...

	select {} // run indefinitely
}
