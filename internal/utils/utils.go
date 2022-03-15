package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Logger struct {
	Log *log.Logger
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

func (l Logger) LogInfo(input string, value string, data interface{}) {
	info := l.Log.WithFields(log.Fields{input: value, "data": data})
	l.Log.Info(info)
}

func (l Logger) LogError(input string, value string, err error) {
	info := l.Log.WithFields(log.Fields{input: value, "error": err.Error()})
	l.Log.Error(info)
}
