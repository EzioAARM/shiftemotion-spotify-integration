package awsFunctions

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
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
