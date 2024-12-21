package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jasonhilder/personal_website/internal/utils"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type SpotifyInfoFailed struct {
    ApiFailed bool
}

// InitSpotify refreshes the access token using the refresh token
func InitSpotify(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

        if(isTokenExpired()) {
            clientID := os.Getenv("SPT_CLIENT_ID")
            clientSecret := os.Getenv("SPT_CLIENT_SECRET")
            refresh_tkn := os.Getenv("SPT_REFRESH_TOKEN")

            url := "https://accounts.spotify.com/api/token"

            // create the x-www-form-urlencoded request body
            data := "grant_type=refresh_token&refresh_token="+refresh_tkn

            // Create the HTTP POST request
            req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
            if err != nil {
                log.Println("Error creating request:", err)
                return
            }

            // Set the necessary headers
            req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
            req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))

            // Send the request using http.Client
            client := &http.Client{}
            resp, err := client.Do(req)
            if err != nil {
                log.Println("Error sending request:", err)
                return
            }
            defer resp.Body.Close()

            // Read the response body for debugging
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                log.Println("Error reading response body:", err)
                return
            }

            log.Printf("%s\n",body)

            // Parse the response body into the TokenResponse struct
            var tokenResponse TokenResponse
            if err := json.Unmarshal(body, &tokenResponse); err != nil {
                log.Println("Error decoding response JSON:", err)
                return
            }
            tokenReceivedTime := time.Now().UnixMilli() 
            tokenExpiresIn := int64(tokenResponse.ExpiresIn * 1000)
            tokenExpiryTime := strconv.Itoa(int(tokenReceivedTime + tokenExpiresIn))

            // Update the environment variable in the current process
            os.Setenv("SPT_TOKEN_EXPIRY", tokenExpiryTime)
            log.Println("Environment variable SPT_TOKEN_EXPIRY updated in current process.")
            os.Setenv("SPT_ACCESS_TOKEN", tokenResponse.AccessToken) 
            log.Println("Environment variable SPT_ACCESS_TOKEN updated in current process.")
        }

		// Call the next handler
		next(w, r)
	}
}

// todo return error....
func isTokenExpired() bool {
    expStamp := os.Getenv("SPT_TOKEN_EXPIRY")
    if expStamp == "" {
		log.Println("Failed to get SPT_TOKEN_EXPIRY:")
    }

    nowStamp := time.Now().UnixMilli()

	// Convert the environment variable to int64
	intExpStamp, err := strconv.ParseInt(expStamp, 10, 64)
	if err != nil {
		log.Printf("Error converting MY_ENV_VAR to int64: %v\n", err)
        return false
	}

    if(intExpStamp == 0 || nowStamp > intExpStamp) {
        return true
    }

    return false
}

func GetSpotifyInfo(w http.ResponseWriter, r *http.Request) {
    access_token := os.Getenv("SPT_ACCESS_TOKEN")
    url := "https://api.spotify.com/v1/me/player?market=za"

    // Create the HTTP POST request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Error creating request:", err)
        return
    }

    // Set the necessary headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+ access_token)

    // Send the request using http.Client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

    if(resp.Status == "200 OK") {
        // Read the response body for debugging
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Println("Error reading response body:", err)
            return
        }

        // Unmarshal the JSON response into a map
        var response utils.SpotifyResponse
        if err := json.Unmarshal(body, &response); err != nil {
            log.Println("Error decoding response JSON:", err)
            return
        }
        response.ApiFailed = false
        
        RenderPage(w, r, "music.html", response)
    } else {
        i := SpotifyInfoFailed{
            ApiFailed: true,
        }

        RenderPage(w, r, "music.html", i)
    }
}
