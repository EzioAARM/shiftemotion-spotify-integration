package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		var suma = 0
		for i := 0; i <= 150000; i++ {
			suma += i
		}
		c.String(200, "Servicio Levantado utilizando AWS codepipeline! (Ah perroooooooooooooooooooo!) "+strconv.Itoa(suma), nil)

	})

	r.Run(":3000")

}
