package main

import (
	"fmt"
	"gb/internal"
	"net/http"
	"sync/atomic"
	"time"
)

func main() {
	serverConfig, err := internal.ConstrServerCfgFromFile()
	if err != nil {
		panic(err)
	}
	greenBlueConfig, err := internal.ConstrBlueGreenCfgFromFile()
	if err != nil {
		panic(err)
	}
	sharedState := &atomic.Bool{}
	sharedState.Store(greenBlueConfig.StartStateIsGreen)
	mux := http.NewServeMux()
	mux.HandleFunc("/set_state", internal.SpawnSetStateHandler(sharedState))
	mux.HandleFunc("/get_state", internal.SpawnGetStateHandler(sharedState))
	mux.HandleFunc("/", internal.SpawnRedirect(sharedState, greenBlueConfig))
	srv := http.Server{Addr: fmt.Sprintf(":%s", serverConfig.Port), Handler: mux, WriteTimeout: 5 * time.Second}
	srv.ListenAndServe()

}
