package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"blockchain/core"
	"blockchain/utils"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000"}

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

type getblocks struct {
	AddrFrom string
}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := core.NewBlockchain(nodeID)

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func sendVersion(addr string, bc *core.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(version{nodeVersion, bestHeight, nodeAddress})

	request := append(utils.CommandToBytes("version"), payload...)

	sendData(addr, request)
}

func handleConnection(conn net.Conn, bc *core.Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := utils.BytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleVersion(request []byte, bc *core.Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func nodeIsKnown(address string) bool {
	for _, addr := range knownNodes {
		if addr == address {
			return true
		}
	}

	return false
}

func sendGetBlocks(address string) {
	payload := utils.GobEncode(getblocks{nodeAddress})
	request := append(utils.CommandToBytes("getblocks"), payload...)

	sendData(address, request)
}
