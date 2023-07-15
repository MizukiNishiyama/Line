package dao

import (
	"database/sql"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"line/model"
)

type UserDao struct {
	DB *sql.DB
}

// 名前で検索
func (dao *UserDao) FindByName(UserId string) ([]model.User, error) {
	rows, err := dao.DB.Query("SELECT UserId, UserName, UserPassword FROM user WHERE UserId = ?", UserId)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.UserId, &u.UserName, &u.UserPassword); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// ユーザー情報を追加
func (dao *UserDao) Insert(user model.User) error {
	_, err := dao.DB.Exec("INSERT into user VALUES(?, ?, ?)", user.UserId, user.UserName, user.UserPassword)
	return err
}

// サインアップ
func (dao *UserDao) Signup(user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = dao.DB.Exec("INSERT INTO user VALUES (?, ?, ?)", user.UserId, user.UserName, string(hashedPassword))
	if err != nil {
		return err
	}

	// Get all existing users.
	rows, err := dao.DB.Query("SELECT UserId, UserName FROM user WHERE UserId != ?", user.UserId)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create a new room for the new user and each existing user.
	for rows.Next() {
		var existingUserID, existingUserName string
		if err := rows.Scan(&existingUserID, &existingUserName); err != nil {
			return err
		}
		RoomId := ulid.Make().String()
		_, err := dao.DB.Exec(
			"INSERT INTO room (RoomId,UserId1, UserId2, UserName1, UserName2) VALUES (?, ?, ?, ?, ?)",
			RoomId, user.UserId, existingUserID, user.UserName, existingUserName)
		if err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// ログイン
func (dao *UserDao) Login(user *model.User) (*model.User, error) {
	result := dao.DB.QueryRow("SELECT UserId, UserName, UserPassword FROM user WHERE UserName = ?", user.UserName)

	storedUser := &model.User{}
	err := result.Scan(&storedUser.UserId, &storedUser.UserName, &storedUser.UserPassword)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.UserPassword), []byte(user.UserPassword))
	if err != nil {
		return nil, err
	}

	return storedUser, nil
}
