package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TCPMessage struct {
	action   int    `json:"action"`
	login    string `json:"login"`
	roomcode string `json:"roomcode"`
	data     string `json:"data"`
}
type UserRoom struct {
	Room     string
	RoomName string
	Role     int
	Image1   string
	Image2   string
	Image3   string
	Events   []Event
}

type Event struct {
	Event  string
	Author string
	Flag   bool
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
	switch msg.action {
	case 0:
		sqlStatement := "SELECT room, status FROM userRoom WHERE user = $1"

		rows, err := db.Query(sqlStatement, msg.login)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		userRooms := []UserRoom{}

		var room string
		var role int
		for rows.Next() {
			userRoom := UserRoom{}
			err = rows.Scan(&room, &role)
			if err != nil {
				panic(err)
			}
			userRoom.Room = room
			userRoom.Role = role

			sqlStatement = "SELECT roomName, image1, image2, image3 FROM roomData WHERE room = $1"
			row := db.QueryRow(sqlStatement, room)
			err := row.Scan(&userRoom.RoomName, &userRoom.Image1, &userRoom.Image2, &userRoom.Image3)
			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("No rows found")
				} else {
					panic(err)
				}
			}
			sqlStatement = `SELECT event, creator, status FROM roomData WHERE room = $1`
			rows, err := db.Query(sqlStatement, room)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			for rows.Next() {
				event := Event{}
				err = rows.Scan(&event.Event, &event.Author, &event.Flag)
				if err != nil {
					panic(err)
				}
				userRoom.Events = append(userRoom.Events, event)

			}
			err = rows.Err()
			if err != nil {
				panic(err)
			}
			userRooms = append(userRooms, userRoom)
		}

		err = rows.Err()
		if err != nil {
			panic(err)
		}
		fmt.Println(userRooms)
		answer, _ := json.Marshal(userRooms)
		c.Write(answer)
	}

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
