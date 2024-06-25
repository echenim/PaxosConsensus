package paxos

type Proposer struct {
	Node
	proposalID int
	value      Value
}

func (p *Proposer) run() {
	p.proposalID++
	prepareMsg := Message{From: p.Id, Type: "prepare", ProposalID: p.proposalID}
	for _, addr := range p.Peers {
		go p.SendMessage(prepareMsg, addr)
	}

	promises := 0
	majority := len(p.Peers)/2 + 1
	for msg := range p.Incoming {
		if msg.Type == "promise" && msg.ProposalID == p.proposalID {
			promises++
			if promises >= majority {
				proposeMsg := Message{From: p.Id, Type: "propose", ProposalID: p.proposalID, Value: p.value}
				for _, addr := range p.Peers {
					go p.SendMessage(proposeMsg, addr)
				}
				break
			}
		}
	}
}
