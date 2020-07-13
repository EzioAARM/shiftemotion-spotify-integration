package awsFunctions

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sony/sonyflake"
	"log"
	"os"
	"strconv"
)

type photo struct {
	ID          string
	UserID      string
	PictureCode string
	Emotion     string
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
	// create random ID
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return err
	}

	item := photo{
		ID:          strconv.FormatUint(id, 10),
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
