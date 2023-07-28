package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
	_ "github.com/oklog/ulid/v2"
	"line/controller"
	"line/dao"
	"line/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var db *sql.DB

func main() {
	db = initDB()

	userDao := &dao.UserDao{DB: db}
	loginUserController := &controller.LoginUserController{LoginUserUseCase: &usecase.LoginUserUseCase{UserDao: userDao}}
	registerUserController := &controller.RegisterUserController{RegisterUserUseCase: &usecase.RegisterUserUseCase{UserDao: userDao}}
	searchUserController := &controller.SearchUserController{SearchUserUseCase: &usecase.SearchUserUseCase{UserDao: userDao}}
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodPost:
			registerUserController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodPost:
			loginUserController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/searchuser", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodGet:
			searchUserController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	messageDao := &dao.MessageDao{DB: db}
	SendMessageController := &controller.SendMessageController{SendMessageUseCase: &usecase.SendMessageUseCase{MessageDao: messageDao}}
	FetchMessageController := &controller.FetchMessageController{FetchMessageUseCase: &usecase.FetchMessageUseCase{MessageDao: messageDao}}

	http.HandleFunc("/sendmessage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodPost:
			SendMessageController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/fetchmessage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodGet:
			FetchMessageController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	roomDao := &dao.RoomDao{DB: db}
	SelectRoomController := &controller.SelectRoomController{SelectRoomUseCase: &usecase.SelectRoomUseCase{RoomDao: roomDao}}
	MakeRoomController := &controller.MakeRoomController{MakeRoomUseCase: &usecase.MakeRoomUseCase{RoomDao: roomDao}}
	http.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodPost:
			SelectRoomController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
	http.HandleFunc("/follow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodPost:
			MakeRoomController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	closeDBWithSysCall()

	//log.Println("Listening...")
	//if err := http.ListenAndServe(":8000", nil); err != nil {
	//	log.Fatal(err)
	//}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initDB() *sql.DB {
	// DB接続のため
	//mysqlUser := os.Getenv("MYSQL_USER")
	//mysqlPwd := os.Getenv("MYSQL_PWD")
	//mysqlHost := os.Getenv("MYSQL_HOST")
	//mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	mysqlUser := "mizuki"
	mysqlPwd := "ramen102"
	mysqlHost := "35.184.81.138"
	mysqlDatabase := "line"

	//mysqlUser := "test_user"
	//mysqlPwd := "password"
	//mysqlHost := "localhost:3306"
	//mysqlDatabase := "test_database"

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)

	_db, err := sql.Open("mysql", connStr)

	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	return _db
}

func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
