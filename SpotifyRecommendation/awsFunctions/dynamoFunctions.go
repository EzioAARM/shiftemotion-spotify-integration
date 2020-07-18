package awsFunctions

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"math/rand"
	"os"
	"time"
)

type photo struct {
	ID          int32  `json:"id"`
	UserID      string `json:"user_id"`
	PictureCode string `json:"picture_code"`
	Emotion     string `json:"emotion"`
}

type refresh struct {
	ID    string `json:"id"`
	Token string `json:"token" `
}

type recomendation struct {
	Id         int32  `json:"id"`
	User       string `json:"user"`
	S3Code     string `json:"s3_code"`
	SongID     string `json:"song_id"`
	SongName   string `json:"song_name"`
	SongArtist string `json:"song_artist"`
}

type OutputTrack struct {
	Name    string `json:"name"`
	Artists string `json:"artists"`
	Id      string `json:"id"`
}

func InsertItem(userId, pictureCode, emotion string) error {
	// Initialize a session
	awsRegion := os.Getenv("REGION")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(awsRegion),
			Credentials: credentials.NewEnvCredentials()},
		Profile: "",
	}))
	// create a new dynamoDB Client
	db := dynamodb.New(sess)
	tableName := "HistorialFotosEmociones"

	item := photo{
		ID:          int32(time.Now().Unix()),
		UserID:      userId,
		PictureCode: pictureCode,
		Emotion:     emotion,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println(err)
		return err
	}

	toInsert := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(toInsert)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func FetchRefresh(email string) (string, error) {
	// Initialize a session
	awsRegion := os.Getenv("REGION")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(awsRegion),
			Credentials: credentials.NewEnvCredentials()},
		Profile: "",
	}))
	// create a new dynamoDB Client
	db := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(email),
			},
		},
		TableName: aws.String("PasswordsTokens"),
	}

	result, err := db.GetItem(input)
	if err != nil {
		return "", err
	}

	var resultToken refresh
	dynamodbattribute.UnmarshalMap(result.Item, &resultToken)

	return resultToken.Token, nil
}

func InsertSongRecommendation(songs []OutputTrack, userId, s3Code string) error {
	awsRegion := os.Getenv("REGION")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(awsRegion),
			Credentials: credentials.NewEnvCredentials()},
		Profile: "",
	}))
	// create a new dynamoDB Client
	db := dynamodb.New(sess)
	tableName := "Recomendaciones"

	// Loop through the input array to insert the data
	position := 0
	for position < len(songs) {
		item := recomendation{
			Id:         int32(time.Now().Unix()) + int32(rand.Intn(100)),
			User:       userId,
			S3Code:     s3Code,
			SongID:     songs[position].Id,
			SongName:   songs[position].Name,
			SongArtist: songs[position].Artists,
		}
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Println(err)
		}

		toInsert := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = db.PutItem(toInsert)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
		}
		if err == nil {
			position++
		} else {
			err = nil
		}
	}

	return nil
}
