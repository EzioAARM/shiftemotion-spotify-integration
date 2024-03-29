package routes

import (
	"../awsFunctions"
	"../songRecommendation"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()
	// Allow CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "accept", "origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "/"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Servicio Levantado", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		fileName, err := awsFunctions.UploadFile(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Logic after the file was successfully uploaded
		fmt.Println(fileName)
		// Logic with Rekognition
		emotion, err := awsFunctions.GetEmotion(fileName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// insert into dynamo db
		userID := c.Request.FormValue("userID")
		err = awsFunctions.InsertItem(userID, fileName, emotion)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Call the Spotify API and return the list of songs
		tracks, err := songRecommendation.FetchSongs(emotion, userID)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		fmt.Println(tracks)
		err = awsFunctions.InsertSongRecommendation(tracks, userID, fileName)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(200, gin.H{
			"tracks":     tracks,
			"statusCode": http.StatusOK,
			"emotion":    emotion,
		})

	})

	r.GET("/fetch", func(c *gin.Context) {
		values := c.Request.URL.Query()
		email := values["email"][0]

		tracks, err := songRecommendation.FetchSongs("SAD", email)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		fmt.Println(tracks)
		err = awsFunctions.InsertSongRecommendation(tracks, "luisfer.marroquin@hotmail.com", "pictures/Simon")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(200, gin.H{
			"tracks":     tracks,
			"statusCode": http.StatusOK,
		})
	})

	return r
}
