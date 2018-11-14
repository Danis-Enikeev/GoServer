package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetUserRoom(t *testing.T) {
	db, _ := sql.Open("sqlite3", "./serverData.db")
	defer db.Close()
	msg := TCPMessage{}
	msg.Action = 0
	msg.Login = "testUser4"
	msg.Roomcode = "-"
	msg.Data = "-"
	answer := getUserRooms(msg, db)
	if len(answer) == 0 {
		t.Errorf("getUserRoom test failed!")
	}
}

func TestGetUserName(t *testing.T) {
	db, _ := sql.Open("sqlite3", "./serverData.db")
	defer db.Close()
	msg := TCPMessage{}
	msg.Action = 1
	msg.Login = "testUser4"
	msg.Roomcode = "-"
	msg.Data = "-"
	answer := getUserName(msg, db)
	if answer.Name != "Name1" {
		t.Errorf("getUserName test failed!")
	}
}

func TestCreateRoom(t *testing.T) {
	db, _ := sql.Open("sqlite3", "./serverData.db")
	defer db.Close()
	msg := TCPMessage{}
	msg.Action = 2
	msg.Login = "testUser4"
	msg.Roomcode = "-"
	msg.Data = "TESTROOM"
	answer := createRoom(msg, db)
	if answer.RoomName != "TESTROOM" {
		t.Errorf("createRoom test failed!")
	}
}
func TestRandStringRunes(t *testing.T) {
	answer := randStringRunes(20)
	if len(answer) != 20 {
		t.Errorf("randStringRunes test failed!")
	}
}
