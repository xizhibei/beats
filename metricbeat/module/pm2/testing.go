package pm2

import "os"

// Helper functions for testing pm2 metricsets.

// GetEnvHost returns the hostname of the pm2 server to use for testing.
// It reads the value from the PM2_HOST environment variable and returns
// 127.0.0.1 if it is not set.
func GetEnvHost() string {
	host := os.Getenv("PM2_HOST")

	if len(host) == 0 {
		host = "127.0.0.1"
	}
	return host
}

// GetEnvPort returns the port of the pm2 server to use for testing.
// It reads the value from the PM2_PORT environment variable and returns
// 9615 if it is not set.
func GetEnvPort() string {
	port := os.Getenv("PM2_PORT")

	if len(port) == 0 {
		port = "9615"
	}
	return port
}
