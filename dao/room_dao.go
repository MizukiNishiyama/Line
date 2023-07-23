package dao

import (
	"database/sql"
	"fmt"
	"github.com/oklog/ulid/v2"
	"line/model"
	"log"
)

type RoomDao struct {
	DB *sql.DB
}

// 指定したRoomを検索
func (dao *RoomDao) SearchRoom(UserId string) ([]model.Room, error) {
	rows, err := dao.DB.Query("SELECT * FROM room WHERE (UserId1 = ? OR UserId2 =?)", UserId, UserId)
	if err != nil {
		return nil, err
	}

	rooms := make([]model.Room, 0)
	for rows.Next() {
		var u model.Room
		if err := rows.Scan(&u.RoomId, &u.UserId1, &u.UserId2, &u.UserName1, &u.UserName2); err != nil {
			return nil, err
		}
		rooms = append(rooms, u)
	}

	return rooms, nil
}

func (dao *RoomDao) MakeRoom(follow model.Follow) error {
	rows, err := dao.DB.Query("SELECT UserId FROM user WHERE UserName = ?", follow.OpponentUserName)
	if err != nil {
		return err
	}
	log.Println(follow.OpponentUserName)
	var OpponentUserId string
	if rows.Next() {
		err := rows.Scan(&OpponentUserId)
		if err != nil {
			return err
		}
		log.Println(OpponentUserId)
	} else {
		return fmt.Errorf("no user found with username %s", follow.OpponentUserName)
	}
	defer rows.Close()
	
	var UserId1, UserId2, UserName1, UserName2 string
	if follow.UserId < OpponentUserId {
		UserId1 = follow.UserId
		UserId2 = OpponentUserId
		UserName1 = follow.UserName
		UserName2 = follow.OpponentUserName
	} else {
		UserId2 = follow.UserId
		UserId1 = OpponentUserId
		UserName2 = follow.UserName
		UserName1 = follow.OpponentUserName
	}

	RoomId := ulid.Make().String()
	_, err = dao.DB.Exec(
		"INSERT INTO room (RoomId,UserId1, UserId2, UserName1, UserName2) VALUES (?, ?, ?, ?, ?)",
		RoomId, UserId1, UserId2, UserName1, UserName2)
	if err != nil {
		return err
	}
	return nil
}
