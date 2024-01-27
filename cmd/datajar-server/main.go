//go:build darwin

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
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

	defaultHandler := func(w http.ResponseWriter, r *http.Request) {
		accessType := "read"
		if r.Method == "POST" {
			accessType = "write"
		}

		whois, remoteIP, err := server.CheckRequestPermissions(lc, r, accessType)
		if err != nil {
			if whois != nil {
				slog.Error("unauthorized", "remoteAddr", r.RemoteAddr, "path", r.URL.Path)
				server.UnauthorizedTotal.WithLabelValues("", "", r.URL.Path).Inc()
			} else {
				slog.Warn("unauthorized", "hostname", whois.Node.Name, "ip", remoteIP, "path", r.URL.Path)
				server.UnauthorizedTotal.WithLabelValues(whois.Node.Name, remoteIP, r.URL.Path).Inc()
			}
			errMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			http.Error(w, errMsg, http.StatusUnauthorized)
			return
		}

		if accessType == "write" {
			slog.Info("write", "hostname", whois.Node.Name, "ip", remoteIP, "path", r.URL.Path)
			server.WritesTotal.WithLabelValues(whois.Node.Name, remoteIP, r.URL.Path).Inc()
			server.HandleWrite(w, r)
		} else {
			slog.Info("read", "hostname", whois.Node.Name, "ip", remoteIP, "path", r.URL.Path)
			server.ReadsTotal.WithLabelValues(whois.Node.Name, remoteIP, r.URL.Path).Inc()
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
	http.Handle("/-/metrics", server.MetricsHandler)
	log.Fatal(http.Serve(ln, nil))
}
