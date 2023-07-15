package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

type Message struct {
	MessageId      string `json:"MessageId"`
	MessageContent string `json:"MessageContent"`
	MessageTime    string `json:"MessageTime"`
	UserId         string `json:"UserId"`
	RoomId         string `json:"RoomId"`
	UserName       string `json:"UserName"`
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

var hub = Hub{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan Message),
}
var db *sql.DB

func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {

			// 必要に応じて、許可するオリジンを指定します
			// 例: return r.Header.Get("Origin") == "http://localhost:3000"
			return true // すべてのオリジンを許可する場合
		},
	}

	// WebSocketの接続を開く
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// handle error
		log.Printf("Failed to upgrade connection: %v", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return // IMPORTANT: Don't forget to return here.
	}

	// 接続終了時に安全に削除
	defer func() {
		delete(hub.clients, ws)
		ws.Close()
	}()

	hub.clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			// handle error
			break
		}

		// 受信したメッセージをブロードキャストする
		hub.broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-hub.broadcast
		for client := range hub.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				// handle error
				client.Close()
				delete(hub.clients, client)
			}
		}
	}
}

func main() {
	db = initDB()

	userDao := &dao.UserDao{DB: db}
	loginUserController := &controller.LoginUserController{LoginUserUseCase: &usecase.LoginUserUseCase{UserDao: userDao}}
	registerUserController := &controller.RegisterUserController{RegisterUserUseCase: &usecase.RegisterUserUseCase{UserDao: userDao}}
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
			SendMessageController.HandleWS(w, r)
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
		case http.MethodGet:
			SelectRoomController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			fmt.Printf(r.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	closeDBWithSysCall()

	log.Println("Listening...")
	//if err := http.ListenAndServe(":6000", nil); err != nil {
	//	log.Fatal(err)
	//}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initDB() *sql.DB {
	// DB接続のための準備
	//mysqlUser := os.Getenv("MYSQL_USER")
	//mysqlPwd := os.Getenv("MYSQL_PWD")
	//mysqlHost := os.Getenv("MYSQL_HOST")
	//mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	//mysqlUser := "uttc"
	//mysqlPwd := "ramen102"
	//mysqlHost := "34.27.193.191:3306"
	//mysqlDatabase := "hackathon"

	mysqlUser := "test_user"
	mysqlPwd := "password"
	mysqlHost := "(localhost:3306)"
	mysqlDatabase := "test_database"

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)

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
