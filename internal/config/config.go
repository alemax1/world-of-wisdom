package config

import (
	"os"
	"strconv"
	"time"
)

type Cfg struct {
	Port                uint16
	ConnectionTimeout   time.Duration
	ConnectionsLimit    int
	ChallengeDifficulty uint8
	ChallengeSize       uint32
	DataFilePath        string
	WorkersCount        int
}

const (
	serverPort          = 9999
	serverConnTimeout   = time.Second * 5
	serverConnsLimit    = 10
	challengeDifficulty = 20
	challengeSize       = 64
	storageDataPath     = "./quotes.json"
	serverWorkersCount  = 10
)

const (
	timeDay = time.Hour * 24
)

func New() *Cfg {
	return &Cfg{
		Port:                uint16(getEnvIntOrDefault("SERVER_PORT", serverPort)),
		ConnectionTimeout:   getDefaultDurationEnv("SERVER_CONNECTION_TIMEOUT", serverConnTimeout),
		ConnectionsLimit:    getEnvIntOrDefault("SERVER_CONNECTIONS_LIMIT", serverConnsLimit),
		ChallengeDifficulty: uint8(getEnvIntOrDefault("CHALLENGE_DIFFICULTY", challengeDifficulty)),
		ChallengeSize:       uint32(getEnvIntOrDefault("CHALLENGE_SIZE", challengeSize)),
		DataFilePath:        getEnvOrDefault("STORAGE_DATA_PATH", storageDataPath),
		WorkersCount:        getEnvIntOrDefault("SERVER_WORKERS_COUNT", serverWorkersCount),
	}
}

func getEnvOrDefault(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}

func getEnvIntOrDefault(key string, def int) int {
	if val, ok := os.LookupEnv(key); ok {
		result, err := strconv.Atoi(val)
		if err != nil {
			return def
		}

		return result
	}

	return def
}

func getDefaultDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if value == "" || !ok {
		return defaultValue
	}

	duration, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return defaultValue
	}

	switch value[len(value)-1] {
	case 'm':
		return time.Duration(duration) * time.Minute
	case 's':
		return time.Duration(duration) * time.Second
	case 'h':
		return time.Duration(duration) * time.Hour
	case 'd':
		return time.Duration(duration) * timeDay
	}

	return defaultValue
}
