package awsFunctions

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// set the region
var awsRegion string

func UploadFile(c *gin.Context) (string, error) {
	awsRegion := os.Getenv("REGION")
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		return "", errors.New("Error parsing the file from the form")
	}
	defer file.Close()

	fmt.Println("**************************")
	fmt.Println(awsRegion)
	// Create the AWS Session
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		log.Println("Error stablishing the session with AWS")
		return "", errors.New("Error stablishing a session with AWS")
	}

	fileName, err := uploadFileToS3(s, file, fileHeader)
	if err != nil {
		log.Println(err)
		return "", errors.New("Error uploading the file to S3")
	}
	fmt.Println("File Successfully Uploaded: ", fileName)

	return fileName, nil
}

func uploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// get the file size and pass the file into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	// get the bucket name
	bucketName := os.Getenv("S3_IMAGE_BUCKET")
	// create a unique file name for the file
	tempFileName := "pictures/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucketName),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("STANDARD"),
	})
	if err != nil {
		return "", err
	}
	return tempFileName, nil
}
