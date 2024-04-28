package utils

import (
	"Omnichannel-CRM/package/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

func init() {
	config.GetConfig()
}

type ResponseData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func RefreshToken() (accessToken string, err error) {
	refreshToken := viper.GetString("OAuth.Refresh_Token")
	clientSecret := viper.GetString("OAuth.Client_Secret")
	clientId := viper.GetString("OAuth.Client_Id")

	reqUrl := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("refresh_token", refreshToken)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Define a struct for the JSON response
	var responseData ResponseData

	// Unmarshal JSON into struct
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", err
	}

	return responseData.AccessToken, nil
}

func NewGmailService() *gmail.Service {
	accessToken := viper.GetString("OAuth.Access_Token")

	accessToken, err := RefreshToken()
	if err != nil {
		accessToken = viper.GetString("OAuth.Access_Token")
	}
	// Create an OAuth2 token source using the access token
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	// Create an OAuth2 HTTP client from the token source
	oauth2Client := oauth2.NewClient(context.Background(), tokenSource)

	// Create a Gmail service using the OAuth2 client
	gmailService, err := gmail.New(oauth2Client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	return gmailService
}

func ParseFromHeader(input string) (name, email string) {
	// Find the position of "<" and ">"
	openBracket := strings.Index(input, "<")
	closeBracket := strings.Index(input, ">")

	// Extract name and email based on the bracket positions
	if openBracket != -1 && closeBracket != -1 {
		name = strings.TrimSpace(input[:openBracket])
		email = strings.TrimSpace(input[openBracket+1 : closeBracket])
	}

	return name, email
}
