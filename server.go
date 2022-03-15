package main

import (
	"github.com/coma-toast/ResumAPI/internal/utils"
	"github.com/gorilla/mux"
)

type API struct {
	conf *utils.Config

	env *Env
}

func (api API) RunAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/", api.LandingHandler)
}
