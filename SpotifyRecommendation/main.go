package main

import (
	"../SpotifyRecommendation/routes"
)

func main() {
	r := routes.RegisterRoutes()

	// Testing purpose

	r.Run(":3000")
}

// Comment this function before pushing and getting in production
