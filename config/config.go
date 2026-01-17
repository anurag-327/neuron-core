package config

import "os"

// PORT: port on which the server will run
var PORT string = "9000"

// Load(): loads the configuration from environment variables and ensure checks are performed
func Load() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	PORT = port
}
