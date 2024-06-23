package P2P

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type NetMessage struct {
	id   string
	IP   string
	port int
}

func sendMessage(peer *Peer) {
	message := NetMessage{
		peer.id,
		peer.IP,
		peer.port,
	}
	messageByte, _ := json.Marshal(message)

	broadcastMessage(messageByte)

}

func broadcastMessage(message []byte) {
	conn, err := net.Dial("udp", "255.255.255.255:8000")
	if err != nil {
		fmt.Println("Error broadcasting message:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending broadcast message:", err)
	}
}

func listenForResponses(peer *Peer, responses chan<- NetMessage) {
	addr := net.UDPAddr{
		Port: 8000,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening for responses:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading UDP message:", err)
			continue
		}
		var response NetMessage
		err = json.Unmarshal(buf[:n], &response)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			continue
		}
		responses <- response
	}
}

func extractPeerFromResponse(response NetMessage) Peer {
	return Peer{
		response.id,
		response.IP,
		response.port,
		[]Peer{},
	}
}

func discoverPeers(peer *Peer) {
	sendMessage(peer)

	responses := make(chan NetMessage)
	go listenForResponses(peer, responses)

	timeout := time.After(10 * time.Second)
	for {
		select {
		case response := <-responses:
			newPeer := extractPeerFromResponse(response)
			known := false
			for _, kp := range peer.known_peers {
				if kp.id == newPeer.id {
					known = true
					break
				}
			}
			if !known {
				peer.known_peers = append(peer.known_peers, newPeer)
			}
		case <-timeout:
			return
		}
	}

}
