package utils

import (
	"fmt"
	"io/fs"
	"strings"
	"time"
)

type Post struct {
    Name string
    PublishDate string
    Url string
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

func GetPost(file fs.DirEntry) Post {
    segments := strings.Split(file.Name(), "_")
    pDate := getTimeFromString(segments[0])
    fName := getNameFromString(strings.Join(segments[1:], " "))
    fUrl := strings.Join(segments[1:], "_")

    p := Post {
        Name: fName,
        PublishDate: pDate,
        Url: fUrl,
    }

    return p
}

