package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/coma-toast/ResumAPI/internal/utils"
	"github.com/coma-toast/ResumAPI/pkg/candidate"
	"github.com/coma-toast/ResumAPI/pkg/nowpush"
	filename "github.com/keepeye/logrus-filename"
	log "github.com/sirupsen/logrus"
)

// One cache to rule them all
var cache = &Cache{}

type Env struct {
	Logger utils.Logger
}

type NowPushInstance interface {
	GetUser() (nowpush.User, error)
	SendMessage(message_type string, note string, link string) (nowpush.MessageResponse, error)
}

type NowPushAPI struct {
	Token  string
	client *nowpush.Client
	env    *Env
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
	env := &Env{Logger: utils.Logger{Log: &log.Logger{}}}
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
	log.SetFormatter(&log.JSONFormatter{})
	if conf.DevMode {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

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

	// Initialize instances
	instances := APIInstances{
		nowPushInstance: &NowPushAPI{
			Token: conf.NowPushAPIKey,
			env:   env,
		},
		candidateDataInstance: &CandidateData{
			env:   env,
			cache: cache,
		},
	}

	api := &API{
		instances: instances,
		conf:      conf,
	}
	go api.RunAPI(env)
	go func() {
		ticker := 0
		for range time.NewTicker(time.Second * time.Duration(300)).C {
			env.Logger.LogInfo("loop", strconv.Itoa(ticker), nil)
			ticker++
		}
	}()

	dontExit := make(chan bool)
	// Waiting for a channel that never comes...
	<-dontExit
}

// * NowPush functions
func (n *NowPushAPI) GetUser() (nowpush.User, error) {
	user := nowpush.User{}
	if n.client == nil {
		n.client = &nowpush.Client{
			Token: n.Token,
		}
	}

	user, err := n.client.GetUser()
	if err != nil {
		n.env.Logger.LogError("error getting client", "nowpush", err)
	}

	return user, err
}
func (n *NowPushAPI) SendMessage(message_type string, note string, link string) (nowpush.MessageResponse, error) {
	if n.client == nil {
		n.client = &nowpush.Client{
			Token: n.Token,
		}
	}

	messageResponse, err := n.client.SendMessage(message_type, note, link)
	if err != nil {
		n.env.Logger.LogError("error sending message", note, err)
	}

	return messageResponse, err
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
