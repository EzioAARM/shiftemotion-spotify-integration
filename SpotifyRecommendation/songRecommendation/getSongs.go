package songRecommendation

import (
	"../awsFunctions"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type tracks struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type tracksArray struct {
	Items []tracks `json:"items"`
}

type token struct {
	Access string `json:"access_token"`
}

type requestInput struct {
	Type  string `url:"grant_type"`
	Token string `url:"refresh_token" `
}

func FetchSongs(emotion, email string) (tracksArray, error) {
	// get the user token
	token, err := retrieveToken(email)
	if err != nil {
		return tracksArray{}, err
	}

	// get a mood
	//spotifyMood := defineMood(emotion)
	topArtists, err := fetchTopTracks(token)
	if err != nil {
		return tracksArray{}, err
	}
	fmt.Print(topArtists)
	// get the key ids from the tracks
	var ids []string
	for _, item := range topArtists.Items {
		ids = append(ids, item.ID)
	}

	fmt.Println(ids)

	return tracksArray{}, nil
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

// Private function to fetch the user's top tracks
// to use them as seeds for recommendations
func fetchTopTracks(token string) (tracksArray, error) {
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks?time_range=short_term&limit=10&offset=0", nil)
	if err != nil {
		return tracksArray{}, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(request)
	if err != nil {
		return tracksArray{}, err
	}

	defer response.Body.Close()
	var retrievedTracks tracksArray
	parsedBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return tracksArray{}, err
	}

	err = json.Unmarshal(parsedBody, &retrievedTracks)
	if err != nil {
		return tracksArray{}, nil
	}

	fmt.Println(retrievedTracks)
	return retrievedTracks, nil
}

// Private function to retrieve a new session token for the user for the current operation
func retrieveToken(email string) (string, error) {
	// get the refresh token
	refreshToken, err := awsFunctions.FetchRefresh(email)
	if err != nil {
		return "", err
	}
	appSecret := "Zjk0MGQ2MTk4OGE2NDg0ZmJkY2M5OGE1OTZkNDc5ZWM6OGZiMzA1ZjA3NzIzNGZhMjhmNjI5YThlYjFmMTI4MmQ="
	// Begin the process of retrieving a new access token
	var access token
	data := requestInput{
		Type:  "refresh_token",
		Token: refreshToken,
	}
	opt, _ := query.Values(data)
	// make the request
	requestToken, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(opt.Encode()))
	if err != nil {
		return "", err
	}
	requestToken.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requestToken.Header.Set("Authorization", "Basic "+appSecret)

	client := &http.Client{Timeout: time.Second * 10}

	response, err := client.Do(requestToken)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	parsedBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// lastly, parse the token
	err = json.Unmarshal([]byte(parsedBody), &access)
	if err != nil {
		return "", err
	}

	// everything is okay
	return access.Access, nil
}
