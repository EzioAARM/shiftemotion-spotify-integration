FROM golang:alpine
RUN apk update && \
    apk upgrade && \
    apk add git
RUN mkdir /app
RUN git clone https://github.com/EzioAARM/shiftemotion-spotify-integration.git /app
WORKDIR /app
RUN go get github.com/gin-gonic/gin
EXPOSE 3000
CMD go run main.go . -DFOREGROUND