package paxos

import "fmt"

type Learner struct {
	Node
	acceptances map[int]map[int]Value // Maps proposal ID to a map of acceptor IDs to values
	majority    int
}

func NewLearner(id int, port int, peers map[int]string, majority int) *Learner {
	return &Learner{
		Node: Node{
			Id:       id,
			Port:     port,
			Incoming: make(chan Message, 10),
			Peers:    peers,
		},
		acceptances: make(map[int]map[int]Value),
		majority:    majority,
	}
}

func (l *Learner) Run() {
	for msg := range l.Incoming {
		if msg.Type == "accepted" {
			if _, exists := l.acceptances[msg.ProposalID]; !exists {
				l.acceptances[msg.ProposalID] = make(map[int]Value)
			}
			l.acceptances[msg.ProposalID][msg.From] = msg.Value

			// Check if we have reached a majority
			if len(l.acceptances[msg.ProposalID]) >= l.majority {
				fmt.Printf("Consensus reached on Proposal %d with Value: %v\n", msg.ProposalID, msg.Value)
			}
		}
	}
}
