package main

import (
	"fmt"
	"os"
	"bufio"
	"net"
	"strconv"
	"encoding/gob"
	"bytes"
	"log"
	hardwaresim "github.com/ntu-brizo/hardwaresim"
	brizochain "github.com/ntu-brizo/brizochain"
)

// goroutine1 (client)
// 1. handle user input
// 2. Encrypt & Sign
// 3. Upload to blockchain
// 4. send hash to the other

// goroutine2 (server)
// 1. handle the other's msg
// 2. Load data from blockchain
// 3. Verify & Decrypt
// 4. Print on screen

type Block1 struct {
    E   []byte      
	K []byte 
	H []byte
	P []byte
	S []byte 
}

func main(){
	listener, err := net.Listen("tcp", ":8081")
    if err != nil {
        log.Fatal("tcp server listener error:", err)
	}
	
	for {
        conn2, err := listener.Accept()
        if err != nil {
            log.Fatal("tcp server accept error", err)
        }

        go handleConnection(conn2)
    }
}


func handleConnection(conn net.Conn) {
	hashString, _ := bufio.NewReader(conn).ReadString('\n')
	log.Println(hashString)

	blockchain, _ := brizochain.NewBrizoChain()
    msgString, _ := blockchain.ReadDataFromHashDict(hashString)

	msg := DeserializeBlock(hardwaresim.HashStringToHash(msgString))
	fmt.Println("================================")
    fmt.Printf("Encrypt content: %x\n",msg.E)
    fmt.Printf("AES key: %x\n",msg.K)
    fmt.Printf("Pub key: %x\n",msg.P)
    fmt.Printf("Hash: %x\n",msg.H)
    fmt.Printf("Sig: %x\n",msg.S)
    fmt.Println("================================")

	valid, _ := hardwaresim.Verify(hashString, msg.S, msg.P)
	fmt.Println(strconv.FormatBool(valid))
	
	content, _ := hardwaresim.Decrypt(msg.K, msg.E)
	contentString := string(content[:])

	clientAddr := conn.RemoteAddr().String()
	response := fmt.Sprintf(contentString + " from " + clientAddr + "\n")
	log.Println(response)
		
    handleConnection(conn)
}

// Serialize serializes the block
func (b *Block1) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block1 {
	var block Block1  

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}