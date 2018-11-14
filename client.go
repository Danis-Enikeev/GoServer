package main

import (
	"encoding/json"
	"net"
	"os"
)

type TCPMessage struct {
	Action   int    `json:"action"`
	Login    string `json:"login"`
	Roomcode string `json:"roomcode"`
	Data     string `json:"data"`
}

func main() {
	message, _ := json.Marshal(TCPMessage{Action: 0, Login: "testUser4", Roomcode: "-", Data: "-"})
	strEcho := message
	servAddr := "localhost:1488"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write(message)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println("write to server = ", strEcho)

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println("reply from server=", string(reply))

	conn.Close()
}
