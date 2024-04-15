// A program to find delvers in our emails and delete their emails so we don't have to read horrifying messages from delvers
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Follow the instructions below:")
	fmt.Printf("1. Go to the following link in your browser"+"\n%v\n", authURL)
	fmt.Println("2.After clicking the link it will send you to your OAuth 2.0 Playground. Click 'Step 1 Select & authorize'")
	fmt.Println("3. Copy and paste this link 'https://mail.google.com/' on the input box provided after clicking 'Step 1 Select & authorize'. Do not click 'Authorize APIs' yet.")
	fmt.Println("4.Click the settings icon on your top right.")
	fmt.Println("5.Click the checkbox 'Use your own OAuth credentials' at the end of 'OAuth 2.0 configuration' form.")
	fmt.Println("6.You will need your 'OAuth Client ID' and 'OAuth Client secret'. Visit your Google APIs Console https://console.cloud.google.com/apis/dashboard?project=delvedetector")
	fmt.Println("7.Click on Credentials on the left sidebar. Under OAuth 2.0 Client IDs click on your project")
	fmt.Println("8.Copy and Paste your 'Client ID' and 'Client secret' to your 'OAuth 2.0 Playground'")
	fmt.Println("9.'Close' the 'OAuth 2.0 configuration' form. Click 'Authorize APIs' next your authorization link")
	fmt.Println("10.You will be redirected to allow for authorization.")
	fmt.Println("11.After approving for authorization. You will get an authorization code. Copy it and paste here")
	fmt.Printf("authorization code:")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	fmt.Println("Checking for delvers...")
	r, _ := srv.Users.Messages.List(user).Do()
	var delvers int
	for _, l := range r.Messages {
		msg, _ := srv.Users.Messages.Get(user, l.Id).Do()
		for _, part := range msg.Payload.Parts {
			if part.MimeType == "text/plain" {
				body, _ := base64.StdEncoding.DecodeString(part.Body.Data)
				body_string := string(body)
				body_string = strings.ToLower(body_string)
				// Let's find that delve among our emails
				searchTerm := "delve"
				if strings.Contains(body_string, searchTerm) {
					deleteEmail(ctx, srv, l.Id)
					delvers++

				}

			}

		}

	}
	fmt.Printf("%v delvers found", delvers)
}

func deleteEmail(ctx context.Context, gmailService *gmail.Service, messageID string) {
	_, err := gmailService.Users.Messages.Trash(
		"me",
		messageID,
	).Do()
	if err != nil {
		fmt.Println(err)
	}
}
