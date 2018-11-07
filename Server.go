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
	Action   int    `json:"action"`
	Login    string `json:"login"`
	Roomcode string `json:"roomcode"`
	Data     string `json:"data"`
}
type UserRoom struct {
	Room     string  `json:"room"`
	RoomName string  `json:"roomName"`
	Role     int     `json:"role"`
	Image1   string  `json:"image1"`
	Image2   string  `json:"image2"`
	Image3   string  `json:"image3"`
	Events   []Event `json:"events"`
}

type Event struct {
	Event  string `json:"event"`
	Author string `json:"author"`
	Flag   bool   `json:"flag"`
}

func handleServerConnection(c net.Conn, db *sql.DB) {

	d := json.NewDecoder(c)
	var msg TCPMessage
	err := d.Decode(&msg)
	if err != nil {
		c.Write([]byte("Wrong input data format"))
		c.Close()
		return
	}
	fmt.Println(msg, err)
	switch msg.Action {
	case 0:
		sqlStatement := "SELECT room, status FROM userRooms WHERE user = $1"

		rows, err := db.Query(sqlStatement, msg.Login)
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
			sqlStatement = `SELECT event, creator, status FROM roomEvents WHERE room = $1`
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
	db, err := sql.Open("sqlite3", "./serverData.db")
	defer db.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("connection established")
		go handleServerConnection(c, db)
	}
}
