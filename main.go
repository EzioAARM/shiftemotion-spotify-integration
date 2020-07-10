package main 

import "github.com/gin-gonic/gin" 

func main() {
	r := gin.Default() 
	
	r.GET("/", func(c *gin.Context){
		c.String(200, "Servicio Levantado implementando cicd!", nil)
	})

	r.Run(":3000") 

}



