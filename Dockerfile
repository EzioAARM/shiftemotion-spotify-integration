FROM golang:alpine
RUN apk update && \
    apk upgrade && \
    apk add git
RUN mkdir /app
RUN git clone https://github.com/EzioAARM/shiftemotion-spotify-integration.git /app
WORKDIR /app
RUN go get github.com/gin-gonic/gin
RUN go get github.com/google/go-querystring/query
RUN go get github.com/gin-contrib/cors
RUN go get github.com/gin-gonic/gin
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/credentials
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/s3
RUN go get github.com/globalsign/mgo/bson
RUN go get github.com/aws/aws-sdk-go/service/dynamodb
RUN go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
RUN go get github.com/sony/sonyflake
RUN go get github.com/aws/aws-sdk-go/aws/awserr
RUN go get github.com/aws/aws-sdk-go/service/rekognition
EXPOSE 3000
CMD go run SpotifyRecommendation/main.go . -DFOREGROUND