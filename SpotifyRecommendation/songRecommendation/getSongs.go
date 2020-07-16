package songRecommendation

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zmb3/spotify"
	"net/http"
	"time"
)

func FetchSongs(emotion, token string) ([]string, error) {
	// get a mood
	//spotifyMood := defineMood(emotion)
	topArtists, err := fetchTopArtist(token)
	if err != nil {
		return nil, err
	}
	fmt.Print(topArtists)
	return nil, nil
}

// Function to convert an emotion from rekognition to a float value spotify can use
// Emotions Valid Values: HAPPY | SAD | ANGRY | CONFUSED | DISGUSTED | SURPRISED | CALM | UNKNOWN | FEAR
func defineMood(mood string) float64 {
	switch mood {
	case "DISGUSTED":
		return 0.0
	case "ANGRY":
		return 0.1
	case "SAD":
		return 0.2
	case "FEAR":
		return 0.3
	case "CONFUSED":
		return 0.4
	case "CALM":
		return 0.7
	case "HAPPY":
		return 0.8
	case "SURPRISED":
		return 0.9
	default:
		return 1.0
	}
}

func fetchTopArtist(token string) ([]spotify.FullArtist, error) {
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/artist?time_range=short_term&limit=10&offset=0", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.New("Error getting response " + err.Error())
	}
	defer response.Body.Close()
	// parse the request
	decoder := json.NewDecoder(response.Body)
	// export the json
	var topArtists []spotify.FullArtist
	err = decoder.Decode(topArtists)

	return topArtists, nil
}
