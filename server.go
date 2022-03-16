package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/coma-toast/ResumAPI/internal/utils"
	"github.com/coma-toast/ResumAPI/pkg/candidate"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type API struct {
	instances APIInstances
	conf      *utils.Config
	env       *Env
}

type APIInstances struct {
	nowPushInstance       NowPushInstance
	candidateDataInstance CandidateDataInstance
}
type JSONResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	Data  string `json:"data,omitempty"`
}

func (api API) RunAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/", api.LandingHandler).Methods(http.MethodGet)
	r.HandleFunc("/", api.AddCandidateHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}/{section}", api.CandidateHandler)
	r.HandleFunc("/{id}", api.SetCandidateHandler).Methods(http.MethodPost)
	r.HandleFunc("/ping", api.PingHandler)
	r.HandleFunc("/reset", api.ResetHandler)
	r.Use(api.notificationMiddleware)
	api.env.Logger.LogError("", "", "api handler failed", http.ListenAndServe(fmt.Sprintf(":%s", api.conf.Port), r))
}

func (api *API) notificationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := api.getIP(r)
		if err != nil {
			api.env.Logger.LogError("", "", "unable to get IP", err)
		}

		message := fmt.Sprintf("API called: %s from %s", r.URL.Path, ip)
		api.env.Logger.LogInfo("API called", r.URL.Path, "", nil)
		api.instances.nowPushInstance.SendMessage("nowpush_note", message, "")

		next.ServeHTTP(w, r)
	})
}

func (api *API) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func (api *API) respondWithError(w http.ResponseWriter, code int, message string) {
	api.respondWithJSON(w, code, JSONResponse{Error: message, OK: false})
}

// PingHandler is just a quick test to ensure api calls are working.
func (api *API) PingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Debug("Request sent to /api/ping")
	api.instances.nowPushInstance.SendMessage("nowpush_note", "PING handler hit", "")

	w.Write([]byte("Pong\n"))
}

func (api *API) LandingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	api.instances.nowPushInstance.SendMessage("nowpush_note", "landing page accessed", "")
	http.Redirect(w, r, api.conf.LandingPage, http.StatusMovedPermanently)
}

func (api *API) CandidateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.env.Logger.LogError("", "", "error getting candidate ID", err)
		api.respondWithError(w, http.StatusBadRequest, "No candidates found.")
		return
	}
	candidate := api.instances.candidateDataInstance.GetCandidateByID(id)
	var data interface{}
	switch section := vars["section"]; section {
	case "contact":
		data = candidate.Contact
	case "experience":
		data = candidate.Experience
	case "projects":
		data = candidate.Projects
	case "dev-env":
		data = candidate.DevEnvs
	case "hobbies":
		data = candidate.Hobbies
	}

	api.respondWithJSON(w, 200, data)
}

func (api *API) SetCandidateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.env.Logger.LogError("", "", "error converting string to int", err)
		api.respondWithError(w, http.StatusBadRequest, "user id error")
		return
	}
	var candidate candidate.Candidate
	err = json.NewDecoder(r.Body).Decode(&candidate)
	if err != nil {
		api.env.Logger.LogError("", "", "error decoding json", err)
		api.respondWithError(w, http.StatusBadRequest, "error setting data")
		return
	}
	err = api.instances.candidateDataInstance.SetCandidate(id, candidate)
	if err != nil {
		api.respondWithError(w, http.StatusBadRequest, "error adding user")
		return
	}
	api.respondWithJSON(w, http.StatusOK, fmt.Sprintf("user %d updated successfully", id))
}

func (api *API) AddCandidateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var candidate candidate.Candidate
	err := json.NewDecoder(r.Body).Decode(&candidate)
	if err != nil {
		api.env.Logger.LogError("", "", "error decoding json", err)
		api.respondWithError(w, http.StatusBadRequest, "error setting data")
		return
	}
	id, err := api.instances.candidateDataInstance.AddCandidate(candidate)
	if err != nil {
		api.respondWithError(w, http.StatusBadRequest, "error adding user")
		return
	}
	api.respondWithJSON(w, http.StatusOK, fmt.Sprintf("user %d added successfully", id))
}

func (api *API) ResetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	candidate := candidate.Candidate{
		// * Put your custom data here - you can use personalData.go to build and object to copy and paste.
		// ! do not commit any changes added here

	}
	id, err := api.instances.candidateDataInstance.AddCandidate(candidate)
	if err != nil {
		api.respondWithError(w, http.StatusBadRequest, "error adding user")
		return
	}
	api.respondWithJSON(w, http.StatusOK, fmt.Sprintf("user %d added successfully", id))
}

func (api *API) getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}
