package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/josh/datajar-server/internal/acl"
	"github.com/josh/datajar-server/internal/datajar/sqlite"

	"tailscale.com/client/tailscale"
	"tailscale.com/tailcfg"
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

	log.Fatal(http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := checkPermissions(lc, r)
		if err != nil {
			errMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		store, err := sqlite.FetchStore()
		if err != nil {
			log.Fatal(err)
		}

		target := acl.GetPath(store, r.URL.Path)

		jsonData, err := json.Marshal(target)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
	})))
}

const peerCapName = "github.com/josh/datajar-server"

func checkPermissions(localClient *tailscale.LocalClient, r *http.Request) error {
	whois, err := localClient.WhoIs(r.Context(), r.RemoteAddr)
	if err != nil {
		return err
	}

	caps, err := tailcfg.UnmarshalCapJSON[acl.Capabilities](whois.CapMap, peerCapName)
	if err != nil {
		return err
	}

	if !acl.CanReadPath(r.URL.Path, caps) {
		return errors.New("unauthorized")
	}

	return nil
}
