package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cp-hackathon/twitch"
	"github.com/cp-hackathon/websocket"
	"github.com/joho/godotenv"
)

func verifySignature(signature string, id string, timestamp string, body []byte) bool {
	message := id + timestamp + string(body)
	h := hmac.New(sha256.New, []byte(os.Getenv("WH_SECRET")))
	h.Write([]byte(message))
	sha := "sha256=" + hex.EncodeToString(h.Sum(nil))
	return signature == sha
}

func handleEvent(pool *websocket.Pool, event twitch.Event) {
	if event.Reward.Title == "Hydrate" {
		log.Println("Hydrate redeemed")
		message := websocket.BroadcastMessage{Event: "hydrate"}
		pool.Broadcast <- message
	}
}

func handleWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := websocket.Upgrade(w, r)

	if err != nil {
		log.Fatal("Error upgrading connection", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client

	log.Println("Client connected to websocket")

	client.Read()
}

func handleNotifcation(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Failed to parse Notification request body", err)
	}

	if !verifySignature(
		r.Header.Get("Twitch-Eventsub-Message-Signature"),
		r.Header.Get("Twitch-Eventsub-Message-Id"),
		r.Header.Get("Twitch-Eventsub-Message-Timestamp"),
		body,
	) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println("Invalid notification signature")
		return
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		var challenge twitch.Challenge

		if err := json.Unmarshal(body, &challenge); err != nil {
			log.Println("Error parsing challenge json", err)
			return
		}

		w.Write([]byte(challenge.Challenge))
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "notification" {
		var notification twitch.Notification

		if err := json.Unmarshal(body, &notification); err != nil {
			log.Println("Error parsing challenge json", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		handleEvent(pool, notification.Event)
	}
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWs(pool, w, r)
	})
	http.HandleFunc("/notification", func(w http.ResponseWriter, r *http.Request) {
		handleNotifcation(pool, w, r)
	})
}

func setupSubscription() {
	token := twitch.GetAppAccessToken(os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))
	twitch.CreateChannelPointsSubscription(
		os.Getenv("TWITCH_BROADCASTER_ID"),
		os.Getenv("TWITCH_CLIENT_ID"),
		token,
		os.Getenv("WH_SECRET"),
		os.Getenv("WH_CALLBACK_URL"),
	)

}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	setupRoutes()
	setupSubscription()

	fmt.Println("Starting server at port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
