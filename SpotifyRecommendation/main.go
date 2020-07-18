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
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA5WKCWMCB5CIXYHRE")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("S3_IMAGE_BUCKET", "shiftemotion-pictures")
	os.Setenv("REGION", "us-east-1")
}
