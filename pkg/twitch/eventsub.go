package twitch

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Condition represents condition data for subscription
type Condition struct {
	BroadCasterID string `json:"broadcaster_user_id"`
	RewardID      string `json:"reward_id"`
}

// Transport represents transport data for subscription
type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

// Subscription represents subscription data
type Subscription struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
	CreatedAt string    `json:"created_at"`
	Cost      int       `json:"cost"`
}

// Challenge represents challenge data from Twitch when sending a subscription request
type Challenge struct {
	Challenge    string       `json:"challenge"`
	Subscription Subscription `json:"subscription"`
}

// Reward represents redeemed channel points reward data
type Reward struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Prompt string `json:"prompt"`
	Cost   int    `json:"cost"`
}

// Event represents event data from EventSub
type Event struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	ID                   string `json:"id"`
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	UserInput            string `json:"user_input"`
	Status               string `json:"status"`
	RedeemedAt           string `json:"redeemed_at"`
	Reward               Reward `json:"reward"`
}

// Notification represents a notification from EventSub
type Notification struct {
	Event        Event        `json:"event"`
	Subscription Subscription `json:"subscription"`
}

// CreateChannelPointsSubscription creates a Twitch EventSub subscription for channel points
func CreateChannelPointsSubscription(broadcasterID string, clientID string, token string, whSecret string, whCallbackURL string) {
	subscription := Subscription{
		Type:    "channel.channel_points_custom_reward_redemption.add",
		Version: "1",
		Condition: Condition{
			BroadCasterID: broadcasterID,
		},
		Transport: Transport{
			Method:   "webhook",
			Callback: whCallbackURL + "/notification",
			Secret:   whSecret,
		},
	}

	postBody, err := json.Marshal(subscription)

	if err != nil {
		log.Fatal("Failed to marshal subscription payload", err)
	}

	req, err := http.NewRequest("POST", BaseHelixURL+"eventsub/subscriptions", bytes.NewBuffer(postBody))

	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	if _, err := client.Do(req); err != nil {
		log.Println("Failed to send subscription request", err)
	}
}
