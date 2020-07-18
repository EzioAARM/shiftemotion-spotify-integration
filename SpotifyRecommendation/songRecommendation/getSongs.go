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

type recommendedTracks struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Artist []struct {
		Name string `json:"name"`
	} `json:"artists"`
}

type tracksArray struct {
	Items []tracks `json:"items"`
}

type recommendations struct {
	Items []recommendedTracks `json:"tracks"`
}

/*
type outputTrack struct {
	Name 		string	`json:"name"`
	Artists 	string 	`json:"artists"`
	Id			string 	`json:"id"`
}*/

type token struct {
	Access string `json:"access_token"`
}

type requestInput struct {
	Type  string `url:"grant_type"`
	Token string `url:"refresh_token" `
}

func FetchSongs(emotion, email string) ([]awsFunctions.OutputTrack, error) {
	// get the user token
	token, err := retrieveToken(email)
	if err != nil {
		return nil, err
	}

	// get a mood
	spotifyMood := defineMood(emotion)
	topArtists, err := fetchTopTracks(token)
	if err != nil {
		return nil, err
	}
	// get the key ids from the tracks
	var ids []string
	for _, item := range topArtists.Items {
		ids = append(ids, item.ID)
	}
	// now that I have the Ids, generate recommendations based on them
	retrievedTracks, err := getRecommendations(ids, fmt.Sprintf("%.2f", spotifyMood), token)
	if err != nil {
		return nil, err
	}
	// assemble the output
	var output []awsFunctions.OutputTrack
	for _, item := range retrievedTracks.Items {
		// Get the list of artist
		names := ""
		for _, art := range item.Artist {
			names += art.Name + ", "
		}

		output = append(output, awsFunctions.OutputTrack{
			Name:    item.Name,
			Artists: names,
			Id:      item.ID,
		})
	}

	return output, nil
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
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks?time_range=short_term&limit=2&offset=0", nil)
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

// Private method to retrieve recommended songs from spotify
// depending on the seed songs entered
func getRecommendations(ids []string, mood, token string) (recommendations, error) {
	// generate the request to spotify
	apiRoute := "https://api.spotify.com/v1/recommendations?limit=20&market=ES&seed_genres=rock%2Cpop%2Creggaeton&seed_tracks=" + strings.Trim(ids[0], " ") + "%2C" + ids[1] + "&target_valence=" + mood
	request, err := http.NewRequest("GET", apiRoute, nil)
	if err != nil {
		return recommendations{}, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(request)
	if err != nil {
		return recommendations{}, err
	}

	defer response.Body.Close()
	var retrievedTracks recommendations
	parsedBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return recommendations{}, err
	}

	err = json.Unmarshal(parsedBody, &retrievedTracks)
	if err != nil {
		return recommendations{}, nil
	}

	fmt.Println("**********************************************************")
	fmt.Println(retrievedTracks)
	return retrievedTracks, nil
}
