package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type CurState struct {
	IsGreen bool `json:"IsGreen"`
}

func SpawnGetStateHandler(isGreenSharedState *atomic.Bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		isGreenNow := isGreenSharedState.Load()
		curState := CurState{IsGreen: isGreenNow}
		result, err := json.Marshal(&curState)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal error", http.StatusBadRequest)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&result)
		}
	}
}
func SpawnSetStateHandler(isGreenSharedState *atomic.Bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var isGreenNew CurState
		err := json.NewDecoder(r.Body).Decode(&isGreenNew)
		if err != nil {
			log.Println(err.Error(), "Internal error, cant set status")
			http.Error(w, "Internal error, cant set status", http.StatusBadRequest)
			return
		}
		log.Printf("Setting is green state to %v", isGreenNew.IsGreen)
		isGreenSharedState.Store(isGreenNew.IsGreen)

		w.WriteHeader(http.StatusNoContent)

	}
}

func SpawnRedirect(isGreenSharedState *atomic.Bool, config BlueGreenConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var redirectTo string
		if isGreenSharedState.Load() {
			redirectTo = config.GreenUrl
		} else {
			redirectTo = config.BlueUrl
		}
		targetUrl, err := url.Parse(redirectTo)
		if err != nil {
			log.Println(redirectTo, err.Error())
			http.Error(w, "Internal", http.StatusBadRequest)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.ServeHTTP(w, r)
	}
}
