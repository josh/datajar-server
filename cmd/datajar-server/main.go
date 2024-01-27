//go:build darwin

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/josh/datajar-server/internal/server"
	"tailscale.com/hostinfo"
	"tailscale.com/tsnet"
)

func main() {
	statedir := flag.String("statedir", "", "Directory to store state")
	hostname := flag.String("hostname", "datajar", "Tailscale node hostname")
	ephemeral := flag.Bool("ephemeral", false, "Register as an Ephemeral node")
	flag.Parse()

	hostinfo.SetApp("datajar")

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

	readHandler := server.CheckRequestPermissionsHandler(lc, "read", server.ReadHandler)
	writeHandler := server.CheckRequestPermissionsHandler(lc, "write", server.WriteHandler)

	http.Handle("/", server.MethodHandler(readHandler, writeHandler))
	http.Handle("/-/healthy", server.HealthyHandler)
	http.Handle("/-/metrics", server.CheckRequestPermissionsHandler(lc, "metrics", server.MetricsHandler))
	log.Fatal(http.Serve(ln, nil))
}
