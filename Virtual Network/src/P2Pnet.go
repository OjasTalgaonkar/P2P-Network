package P2P

import (
	"fmt"
	"net"
)

type Peer struct {
	id          string
	IP          string
	port        int
	known_peers []Peer
}

func createPeer(id string, ip string, port int) *Peer {
	return &Peer{
		id,
		ip,
		port,
		make([]Peer, 0),
	}
}

func is_connected(peer *Peer, knownPeer Peer) bool {
	return peer.IP == knownPeer.IP && peer.port == knownPeer.port
}

func establish_connection(peer *Peer, knownPeer Peer) {
	// Build the address string for the known peer
	address := fmt.Sprintf("%s:%d", knownPeer.IP, knownPeer.port)

	// Attempt to establish a TCP connection to the known peer
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Error connecting to %s: %v\n", address, err)
		return
	}
	defer conn.Close()

	// Successful connection message
	fmt.Printf("Successfully established TCP connection between %s:%d and %s:%d\n",
		peer.IP, peer.port, knownPeer.IP, knownPeer.port)

}
