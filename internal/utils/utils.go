package utils

import (
	"fmt"
	"strings"
	"time"
)

type SpotifyResponse struct {
	IsPlaying  bool  `json:"is_playing"`

	Item struct {

		Album struct {
			Images []struct {
				Url    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			Artists              []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"album"`

		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`

		Name string `json:"name"`

	} `json:"item"`

    ApiFailed bool `json:"-"`
}

type GistId struct {
	Id string
}

func getTimeFromString(inputDate string) string {
    // Parse the input date string into a time.Time object
    parsedDate, err := time.Parse("20060102", inputDate)
    if err != nil {
        fmt.Println("Error parsing date:", err)
    }
    
    // Format the parsed date as desired
    formattedDate := parsedDate.Format("2006-01-02")
    return formattedDate
}

func getNameFromString(name string) string { 
    out := strings.TrimSuffix(name, ".html")

    return out
}

func GetGistId(id string) GistId {
	g := GistId {
		Id: id,
	}

	return g;
}

