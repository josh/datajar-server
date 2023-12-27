package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/josh/datajar-server/internal/server"

	"tailscale.com/tsnet"
)

func main() {
	statedir := flag.String("statedir", "", "Directory to store state")
	hostname := flag.String("hostname", "datajar", "Tailscale node hostname")
	ephemeral := flag.Bool("ephemeral", false, "Register as an Ephemeral node")
	flag.Parse()

	s := &tsnet.Server{
		Dir:       *statedir,
		Hostname:  *hostname,
		Ephemeral: *ephemeral,
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

	defaultHandler := func(w http.ResponseWriter, r *http.Request) {
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
			server.HandleWrite(w, r)
		} else {
			server.HandleRead(w, r)
		}
	}

	if s.Ephemeral {
		c := make(chan os.Signal, 1)
		shutdown := func() {
			<-c
			lc.Logout(context.TODO())
			os.Exit(1)
		}
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go shutdown()
	}

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/-/healthy", server.HandleHealthy)
	log.Fatal(http.Serve(ln, nil))
}
