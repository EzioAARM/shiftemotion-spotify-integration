package main

import (
	"../SpotifyRecommendation/routes"
	"os"
)

func main() {
	r := routes.RegisterRoutes()

	// Testing purpose
	setEnvVariables()

	r.Run(":3000")
}

// Comment this function before pushing and getting in production
func setEnvVariables() {
	os.Setenv("GIN_MODE", "release")
}
