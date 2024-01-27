//go:build darwin

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/josh/datajar-server/internal/datajar/scriptingbridge"
	"github.com/josh/datajar-server/internal/datajar/sqlite"
	"tailscale.com/client/tailscale/apitype"
)

var healthError error

var ReadHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	whois := r.Context().Value(WhoisKey).(*apitype.WhoIsResponse)
	remoteIP := r.Context().Value(RemoteIPKey).(string)

	slog.Info("read", "hostname", whois.Node.Name, "ip", remoteIP, "path", r.URL.Path)
	ReadsTotal.WithLabelValues(whois.Node.Name, remoteIP, r.URL.Path).Inc()

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
})

var WriteHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	whois := r.Context().Value(WhoisKey).(*apitype.WhoIsResponse)
	remoteIP := r.Context().Value(RemoteIPKey).(string)

	slog.Info("write", "hostname", whois.Node.Name, "ip", remoteIP, "path", r.URL.Path)
	WritesTotal.WithLabelValues(whois.Node.Name, remoteIP, r.URL.Path).Inc()

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
})

var HealthyHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if healthError != nil {
		http.Error(w, healthError.Error(), 500)
	} else {
		fmt.Fprintf(w, "OK")
	}
})

func MethodHandler(getHandle http.Handler, postHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getHandle.ServeHTTP(w, r)
		} else if r.Method == "POST" {
			postHandler.ServeHTTP(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}
