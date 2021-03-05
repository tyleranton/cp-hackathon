package twitch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// AccessToken represents access token data from Twitch
type AccessToken struct {
	Token        string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

// GetAppAccessToken retrieves an app access token from Twitch
func GetAppAccessToken(clientID string, clientSecret string) string {
	v := url.Values{}

	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	v.Set("grant_type", "client_credentials")
	v.Set("scope", "channel:read:redemptions channel:manage:redemptions")

	req, err := http.Post(BaseTokenURL+"?"+v.Encode(), "application/json", nil)

	if err != nil {
		log.Println("Failed to get app access token", err)
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println("Failed to parse request body", err)
	}

	var accessToken AccessToken

	if err := json.Unmarshal(body, &accessToken); err != nil {
		log.Fatal("Error parsing json", err)
	}

	return accessToken.Token
}
