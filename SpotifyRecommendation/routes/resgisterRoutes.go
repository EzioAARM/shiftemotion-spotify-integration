package routes

import (
	"../awsFunctions"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()
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
		c.JSON(http.StatusOK, gin.H{
			"emotion": emotion,
		})

	})

	return r
}
