package configs

import "os"

func EnvMongoURI() string {
	if value, ok := os.LookupEnv("MONGO_URI"); ok {
		return value
	}
	return "mongodb+srv://akash:akash@nosana-battery-agent.73qrj2z.mongodb.net/"
}

func EnvPromURL() string {
	if value, ok := os.LookupEnv("PROMETHEUS_URL"); ok {
		return value
	}
	return "http://portal.greenkwh.net:9090"
}

func EnvPort() string {
	if value, ok := os.LookupEnv("PORT"); ok {
		return value
	}
	return "8080"
}

func EnvElizaOSURL() string {
	if value, ok := os.LookupEnv("ELIZAOS_URL"); ok {
		return value
	}
	return "http://localhost:3000"
}
