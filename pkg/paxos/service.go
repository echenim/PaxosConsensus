package paxos

import "sync"

var ConfigService = struct {
	sync.RWMutex
	Peers map[int]string
}{
	Peers: make(map[int]string),
}

func UpdatePeers(node *Node) {
	ConfigService.RLock()
	defer ConfigService.RUnlock()
	node.Peers = ConfigService.Peers
}

func AddNode(id int, addr string) {
	ConfigService.Lock()
	defer ConfigService.Unlock()
	ConfigService.Peers[id] = addr
}

func RemoveNode(id int) {
	ConfigService.Lock()
	defer ConfigService.Unlock()
	delete(ConfigService.Peers, id)
}
