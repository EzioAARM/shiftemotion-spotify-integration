package main

import (
	"../SpotifyRecommendation/routes"
	"os"
)

func main() {
	r := routes.RegisterRoutes()

	// Testing purposes
	setEnvVariables()

	r.Run(":3000")
}

// Comment this function before pushing and getting in prod
func setEnvVariables() {
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA5WKCWMCB4COQMP5A")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	os.Setenv("S3_IMAGE_BUCKET", "shift-emotion-pictures")
	os.Setenv("REGION", "us-east-2")
}
