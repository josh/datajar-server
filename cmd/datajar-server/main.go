package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/josh/datajar-server/internal/datajar/scriptingbridge"
	"github.com/josh/datajar-server/internal/datajar/sqlite"
	"github.com/josh/datajar-server/internal/server"

	"tailscale.com/tsnet"
)

func main() {
	dir := filepath.Join("./state")
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	s := &tsnet.Server{
		Dir:       dir,
		Hostname:  "datajar",
		Ephemeral: true,
	}
	defer s.Close()

	ln, err := s.ListenTLS("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	var healthError error

	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessType := "read"
		if r.Method == "POST" {
			accessType = "write"
		}

		err := server.CheckRequestPermissions(lc, r, accessType)
		if err != nil {
			errMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		if accessType == "write" {
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
		} else {
			store, err := sqlite.FetchStore()
			if err != nil {
				healthError = err
				log.Fatal(err)
			}
			healthError = nil

			target, err := server.GetValueByPath(store, r.URL.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			jsonData, err := json.Marshal(target)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			fmt.Fprintf(w, "%s\n", jsonData)
		}
	}))

	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		if healthError != nil {
			http.Error(w, healthError.Error(), 500)
		} else {
			fmt.Fprintf(w, "OK")
		}
	})

	log.Fatal(http.Serve(ln, nil))
}
