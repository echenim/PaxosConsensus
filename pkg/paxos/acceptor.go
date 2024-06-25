package paxos

import (
	"sync"
)


type Acceptor struct {
    Node
    acceptedID  int
    acceptedVal Value
    mu          sync.Mutex
}

func (a *Acceptor) run() {
    for msg := range a.Incoming {
        a.mu.Lock()
        if msg.Type == "prepare" && msg.ProposalID > a.acceptedID {
            promiseMsg := Message{From: a.Id, Type: "promise", ProposalID: msg.ProposalID}
            go a.SendMessage(promiseMsg, a.Peers[msg.From])
        } else if msg.Type == "propose" && msg.ProposalID >= a.acceptedID {
            a.acceptedID = msg.ProposalID
            a.acceptedVal = msg.Value
            acceptedMsg := Message{From: a.Id, Type: "accepted", ProposalID: a.acceptedID, Value: a.acceptedVal}
            go a.SendMessage(acceptedMsg, a.Peers[msg.From])
        }
        a.mu.Unlock()
    }
}
