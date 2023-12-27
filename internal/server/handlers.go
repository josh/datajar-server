package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/josh/datajar-server/internal/datajar/scriptingbridge"
	"github.com/josh/datajar-server/internal/datajar/sqlite"
)

var healthError error

func HandleRead(w http.ResponseWriter, r *http.Request) {
	store, err := sqlite.FetchStore()
	if err != nil {
		healthError = err
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	healthError = nil

	target, err := GetValueByPath(store, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// If target type is string, just output it, otherwise serialize to json
	if _, ok := target.(string); ok {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "%s\n", target)
		return
	} else {
		jsonData, err := json.Marshal(target)
		if err != nil {
			healthError = err
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		healthError = nil
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", jsonData)
	}
}

func HandleWrite(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.URL.Path == "" || r.URL.Path == "/" {
		http.Error(w, "cannot write to root", http.StatusBadRequest)
		return
	}

	err = scriptingbridge.SetStoreValue(r.URL.Path, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleHealthy(w http.ResponseWriter, r *http.Request) {
	if healthError != nil {
		http.Error(w, healthError.Error(), 500)
	} else {
		fmt.Fprintf(w, "OK")
	}
}
