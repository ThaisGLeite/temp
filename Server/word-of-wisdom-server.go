package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var wisdom = []string{
	"The only true wisdom is in knowing you know nothing.",
	"Do not seek to follow in the footsteps of the wise; seek what they sought.",
	"It is the mark of an educated mind to be able to entertain a thought without accepting it.",
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Create a unique server-provided string (challenge) for each connection
	challenge := strconv.Itoa(rand.Int())

	// Send challenge to client
	writer := bufio.NewWriter(conn)
	writer.WriteString(challenge + "\n")
	writer.Flush()

	// Get response (nonce) from client
	reader := bufio.NewReader(conn)
	nonce, err := reader.ReadString('\n')
	nonce = strings.TrimSpace(nonce)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Hashcash was chosen because it's a proven and straightforward PoW system.
	//The client has to do significant work to find a suitable nonce, but it's quick and easy for the server to verify the work.
	//This asymmetry is what makes PoW an effective method to deter DDoS attacks: attackers have to expend significant resources to spam requests,
	//while the server can verify these requests cheaply. The challenge-response protocol helps ensure that the work done by the client is unique per connection and cannot be precomputed ahead of time.

	// Hash challenge and nonce together and check for leading zeros
	h := sha256.New()
	h.Write([]byte(challenge + nonce))
	hash := hex.EncodeToString(h.Sum(nil))
	if !strings.HasPrefix(hash, "00000") {
		fmt.Println("Invalid PoW")
		return
	}

	// Send wisdom
	rnd := time.Now().UnixNano() % int64(len(wisdom))
	wisdomMsg := wisdom[rnd]
	writer.WriteString(wisdomMsg + "\n")
	writer.Flush()
}
