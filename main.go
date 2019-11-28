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

var (
	wsPort        = flag.Int("wsport", 8080, "Websocket server port")
	redisHost     = flag.String("h", "127.0.0.1", "Redis host")
	redisPort     = flag.Int("p", 6379, "Redis port number")
	redisPassword = flag.String("pass", "", "Redis password")
	redisDBIndex  = flag.Int("db", 0, "Redis DB index")
	channelName   = flag.String("chan", "browsers", "Redis channel name")

	redisClient *redis.Client
	clients     = make(map[*websocket.Conn]*redis.PubSub)
	subCount    = 0
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	flag.Parse()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     *redisHost + ":" + strconv.Itoa(*redisPort),
		Password: *redisPassword,
		DB:       *redisDBIndex,
	})

	go pollRedis()
	go pollDeadClients()

	fmt.Println("Starting WS server...")
	http.HandleFunc("/browsers", wsHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(*wsPort), nil)
	logFatalError(err)
}

func pollRedis() {
	for {
		pubSubNum, err := redisClient.PubSubNumSub(*channelName).Result()
		logFatalError(err)

		newSubCount := int(pubSubNum[*channelName])
		if newSubCount != subCount {
			subCount = newSubCount
			broadcast()
		}

		time.Sleep(time.Second)
	}
}

func pollDeadClients() {
	for {
		for wClient, redisSub := range clients {
			if err := wClient.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				removeClient(redisSub, wClient)
			}
		}

		time.Sleep(time.Second)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	logFatalError(err)
	clients[conn] = redisClient.Subscribe(*channelName)
}

func broadcast() {
	for wClient, redisSub := range clients {
		browsers := strconv.Itoa(subCount)
		err := wClient.WriteMessage(websocket.TextMessage, []byte(browsers))
		if err != nil {
			removeClient(redisSub, wClient)
		}
	}
}

func removeClient(redisSub *redis.PubSub, wClient *websocket.Conn) {
	err := redisSub.Unsubscribe(*channelName)
	logError(err)
	wClient.Close()
	delete(clients, wClient)
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func logFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
