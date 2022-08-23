package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/coma-toast/ResumAPI/internal/notifications"
	"github.com/coma-toast/ResumAPI/internal/utils"
	"github.com/coma-toast/ResumAPI/pkg/candidate"
	filename "github.com/keepeye/logrus-filename"
	log "github.com/sirupsen/logrus"
)

// One cache to rule them all
var cache = &Cache{}

type Env struct {
	Logger utils.Logger
}

type NotificationInstance interface {
	SendMessage(title, body string) error
}

type CandidateDataInstance interface {
	GetCandidateByID(id int) candidate.Candidate
	SetCandidate(id int, date candidate.Candidate) error
	AddCandidate(candidate.Candidate) (int, error)
}

type CandidateData struct {
	env   *Env
	cache *Cache
}

func main() {
	// Get config file and flags
	configFileLocation, err := os.UserHomeDir()
	if err != nil {
		utils.HandleErr("location", "Home Directory error", err)
	}
	configFileLocation = fmt.Sprintf("%s/%s", configFileLocation, ".config/")
	configPath := flag.String("conf", configFileLocation, "Path for the config file. Default is $HOME/.config")
	flag.Parse()
	conf := utils.GetConf(*configPath)

	// Logging setup
	if conf.LogJson {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if conf.DevMode {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	// TODO: logging DB

	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	log.AddHook(filenameHook)
	date := time.Now().Format("01-02-2006")
	err = os.MkdirAll(conf.LogFilePath, 0777)
	if err != nil {
		log.WithField("error", err).Warn("Unable to create log directory")
	}
	fileLocation := fmt.Sprintf("%s/pod-log-%s.log", conf.LogFilePath, date)
	file, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err == nil {
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		log.WithField("error", err).Warn("Failed to log to file, using default stderr")
	}

	// Cache setup
	err = cache.Initialize(conf.CachePath)
	if err != nil {
		log.WithField("error", err).Warn("Error initializing cache")
	}
	env := &Env{Logger: utils.Logger{Log: &log.Logger{}}}

	// Initialize instances
	instances := APIInstances{
		notificationInstance: &notifications.NotificationAPI{
			Target: conf.NotifAPITarget,
		},
		candidateDataInstance: &CandidateData{
			env:   env,
			cache: cache,
		},
	}

	api := &API{
		instances: instances,
		conf:      conf,
		env:       env,
	}
	go api.RunAPI()
	go func() {
		ticker := 0
		for range time.NewTicker(time.Second * time.Duration(300)).C {
			env.Logger.LogInfo("loop", strconv.Itoa(ticker), "", nil)
			ticker++
		}
	}()

	dontExit := make(chan bool)
	// Waiting for a channel that never comes...
	<-dontExit
}

func (c *CandidateData) GetCandidateByID(id int) candidate.Candidate {
	if c.cache == nil {
		c.cache = &Cache{}
	}
	return c.cache.GetCandidateByID(id)
}

func (c *CandidateData) SetCandidate(id int, data candidate.Candidate) error {
	if c.cache == nil {
		c.cache = &Cache{}
	}
	return c.cache.SetCandidate(strconv.Itoa(id), data)
}

func (c *CandidateData) AddCandidate(data candidate.Candidate) (int, error) {
	if c.cache == nil {
		c.cache = &Cache{}
	}

	return c.cache.AddCandidate(data)
}
