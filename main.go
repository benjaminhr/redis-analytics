package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

var redisHost = flag.String("h", "127.0.0.1", "Redis host")
var redisPort = flag.Int("p", 6379, "Redis port number")
var redisPassword = flag.String("pass", "", "Redis password")
var redisDBIndex = flag.Int("db", 0, "Redis DB index")
var channelName = flag.String("chan", "browsers", "Redis channel name")

var redisClient = redis.NewClient(&redis.Options{
	Addr:     *redisHost + ":" + strconv.Itoa(*redisPort),
	Password: *redisPassword,
	DB:       *redisDBIndex,
})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var subCount = 0

func main() {
	flag.Parse()

	go func() {
		for {
			pubSubNum, _ := redisClient.PubSubNumSub(*channelName).Result()
			newSubCount := int(pubSubNum[*channelName])
			if newSubCount != subCount {
				subCount = newSubCount
				broadcast()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/browsers", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		clients[conn] = true

		browsers := strconv.Itoa(subCount)
		err = conn.WriteMessage(websocket.TextMessage, []byte(browsers))
		if err != nil {
			log.Printf("Error sending WS message: %s", err)
		}
	})

	fmt.Println("Starting WS server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func broadcast() {
	for client := range clients {
		browsers := strconv.Itoa(subCount)
		err := client.WriteMessage(websocket.TextMessage, []byte(browsers))
		if err != nil {
			log.Printf("Websocket error: %s", err)
			client.Close()
			delete(clients, client)
		}
	}
}
