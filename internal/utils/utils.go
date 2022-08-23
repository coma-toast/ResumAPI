package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"

	"github.com/shirou/gopsutil/host"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	Log *log.Logger
}

type HealthData struct {
	Uptime uint64
}

// HandleErr is to handle general errors with logging
func HandleErr(field string, value string, err error) {
	logger := log.WithFields(log.Fields{"field": field, "value": value})
	logger.Warn(err)
}

func HashPassword(rawPassword string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(rawPassword)))
}

func Contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func GenerateRandomString(length int) string {
	var output strings.Builder
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func (l Logger) LogInfo(key string, value string, message string, data interface{}) {
	l.Log.WithFields(log.Fields{key: value, "data": data}).Info(message)
}

func (l Logger) LogDebug(message string) {
	l.Log.Debug(message)
}

func (l Logger) LogError(key string, value string, message string, err error) {
	l.Log.WithFields(log.Fields{key: value, "error": err.Error()}).Error(message)
}

func GetHealthData() HealthData {
	uptime, _ := host.Uptime()
	return HealthData{
		Uptime: uptime,
	}
}
