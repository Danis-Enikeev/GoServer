package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const MIN = 1
const MAX = 100

type TCPMessage struct {
	action   int    `json:"action"`
	login    string `json:"login"`
	password string `json:"password"`
	roomcode string `json:"roomcode"`
	data     string `json:"data"`
}

func random() int {
	return rand.Intn(MAX-MIN) + MIN
}

/*func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}

		result := strconv.Itoa(random()) + "\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
}*/
func handleServerConnection(c net.Conn, db *sql.DB) {

	// we create a decoder that reads directly from the socket
	d := json.NewDecoder(c)
	var msg TCPMessage
	err := d.Decode(&msg)
	fmt.Println(msg, err)
	sqlStatement := "INSERT INTO logs (roomcode, login, password) VALUES ($1, $2, $3)"
	_, err = db.Exec(sqlStatement, msg.roomcode, msg.login, msg.data)
	c.Close()

}
func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	rand.Seed(time.Now().Unix())
	db, err := sql.Open("sqlite3", "./serverDB.db")
	defer db.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleServerConnection(c, db)
	}
}
