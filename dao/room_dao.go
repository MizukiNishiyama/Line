package dao

import (
	"database/sql"
	"line/model"
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
