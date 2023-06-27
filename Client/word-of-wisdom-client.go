package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Read the challenge from the server
	reader := bufio.NewReader(conn)
	challenge, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	challenge = strings.TrimSpace(challenge)

	// Solve the PoW challenge
	nonce := solvePow(challenge)

	// Send the solution to the server
	writer := bufio.NewWriter(conn)
	writer.WriteString(nonce + "\n")
	writer.Flush()

	// Read and print the wisdom from the server
	wisdom, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(wisdom)
}

func solvePow(challenge string) string {
	for i := 0; ; i++ {
		nonce := fmt.Sprintf("%d", i)
		h := sha256.New()
		h.Write([]byte(challenge + nonce))
		hash := hex.EncodeToString(h.Sum(nil))
		if strings.HasPrefix(hash, "00000") {
			return nonce
		}
	}
}
