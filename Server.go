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
	Room     string     `json:"room"`
	RoomName string     `json:"roomName"`
	Role     int        `json:"role"`
	Image1   string     `json:"image1"`
	Image2   string     `json:"image2"`
	Image3   string     `json:"image3"`
	Events   []Event    `json:"events"`
	Users    []UserData `json:"users"`
}

type Event struct {
	Event  string `json:"event"`
	Author string `json:"author"`
	Flag   bool   `json:"flag"`
}

type UserData struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func getUserRooms(msg TCPMessage, db *sql.DB) []UserRoom {
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

		sqlStatement = `SELECT user FROM userRooms WHERE room = $1`
		rows, err = db.Query(sqlStatement, room)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			tempMessage := TCPMessage{}
			err = rows.Scan(&tempMessage.Login)
			if err != nil {
				panic(err)
			}
			tempUser := getUserName(tempMessage, db)
			userRoom.Users = append(userRoom.Users, tempUser)

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

	return userRooms
}

func getUserName(msg TCPMessage, db *sql.DB) UserData {
	userData := UserData{}
	sqlStatement := "SELECT name, image FROM userData WHERE user = $1"
	row := db.QueryRow(sqlStatement, msg.Login)
	err := row.Scan(&userData.Name, &userData.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
		} else {
			panic(err)
		}
	}

	return userData
}
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func createRoom(msg TCPMessage, db *sql.DB) UserRoom {
	room := UserRoom{}
	for true {
		stmt, err := db.Prepare("INSERT INTO roomData(room, roomName, image1, image2, image3) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			panic(err)
		}
		roomCode := randStringRunes(8)
		_, err = stmt.Exec(roomCode, msg.Data, "-", "-", "-")
		if err == nil {
			room.Room = roomCode
			room.RoomName = msg.Data
			room.Image1 = "-"
			room.Image2 = "-"
			room.Image3 = "-"
			room.Role = 0
			stmt, err := db.Prepare("INSERT INTO userRooms(user, room, status) VALUES (?, ?, ?)")
			if err != nil {
				panic(err)
			}
			_, err = stmt.Exec(msg.Login, room.Room, room.Role)
			if err != nil {
				panic(err)
			}
			break
		}
	}

	return room
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
	answer := []byte("error, wrong action")
	switch msg.Action {
	case 0:
		answer, _ = json.Marshal(getUserRooms(msg, db)) //get all rooms and data
	case 1:
		answer, _ = json.Marshal(getUserName(msg, db)) //get username and profile pic
	case 2:
		answer, _ = json.Marshal(createRoom(msg, db)) //create room

	case 3:
	case 4:
	}
	c.Write(answer)
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
