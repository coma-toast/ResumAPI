package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (api API) RunAPI(env *Env) {
	r := mux.NewRouter()
	r.HandleFunc("/", api.LandingHandler)
	r.HandleFunc("/{id}/{section}", api.CandidateHandler)
	r.HandleFunc("/{id}", api.SetCandidateHandler).Methods(http.MethodPost)
	r.HandleFunc("/ping", api.PingHandler)
	r.HandleFunc("/reset", api.ResetHandler)
	api.env.Logger.LogError("api handler failed", "", http.ListenAndServe(fmt.Sprintf(":%s", api.conf.Port), r))
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
	http.Redirect(w, r, api.conf.LandingPage, 200)
}

func (api *API) CandidateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.env.Logger.LogError("error getting candidate ID", "", err)
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
	var candidate candidate.Candidate
	err := json.NewDecoder(r.Body).Decode(&candidate)
	if err != nil {
		api.env.Logger.LogError("error decoding json", "", err)
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
		// * Put your custom data here
	}
	id, err := api.instances.candidateDataInstance.AddCandidate(candidate)
	if err != nil {
		api.respondWithError(w, http.StatusBadRequest, "error adding user")
		return
	}
	api.respondWithJSON(w, http.StatusOK, fmt.Sprintf("user %d added successfully", id))
}
