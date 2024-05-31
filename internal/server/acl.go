package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"tailscale.com/client/tailscale"
	"tailscale.com/client/tailscale/apitype"
	"tailscale.com/tailcfg"
)

const PeerCapName = "github.com/josh/datajar-server"

type Capabilities struct {
	Read    []string `json:"read"`
	Write   []string `json:"write"`
	Metrics bool     `json:"metrics"`
}

func CanAccessPath(requestPath string, caps []Capabilities, accessType string) bool {
	requestPath = strings.TrimSuffix(requestPath, "/")
	if requestPath == "" {
		requestPath = "/"
	}

	for _, cap := range caps {
		var accessList []string
		if accessType == "read" {
			accessList = cap.Read
		} else if accessType == "write" {
			accessList = cap.Write
		} else if accessType == "metrics" {
			if cap.Metrics {
				return true
			} else {
				continue
			}
		}

		for _, accessPattern := range accessList {
			pathPrefix := "/" + strings.TrimPrefix(strings.TrimSuffix(accessPattern, "*"), "/")
			if strings.HasPrefix(requestPath, pathPrefix) {
				return true
			}
		}
	}
	return false
}

func CheckRequestPermissions(localClient *tailscale.LocalClient, r *http.Request, accessType string) (*apitype.WhoIsResponse, string, error) {
	host, _, spliterr := net.SplitHostPort(r.RemoteAddr)
	whois, whoiserr := localClient.WhoIs(r.Context(), r.RemoteAddr)

	if spliterr != nil {
		return whois, host, spliterr
	} else if whoiserr != nil {
		return whois, host, whoiserr
	}

	if host[0:4] != "100." {
		return whois, host, errors.New("remoteAddr is not a Tailscale IP")
	}

	caps, err := tailcfg.UnmarshalCapJSON[Capabilities](whois.CapMap, PeerCapName)
	if err != nil {
		return whois, host, err
	}

	if !CanAccessPath(r.URL.Path, caps, accessType) {
		return whois, host, errors.New("unauthorized")
	}

	return whois, host, nil
}

type ACLMiddleware struct {
	handler     http.Handler
	localClient *tailscale.LocalClient
	accessType  string
}

type contextKey string

const WhoisKey contextKey = "whois"
const RemoteIPKey contextKey = "remoteIP"

func (m *ACLMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	whois, remoteIP, err := CheckRequestPermissions(m.localClient, r, m.accessType)
	if err != nil {
		if whois != nil {
			slog.Warn("unauthorized", "hostname", whois.Node.Name, "ip", remoteIP, "path", r.URL.Path)
			UnauthorizedTotal.WithLabelValues(whois.Node.Name, remoteIP, r.URL.Path).Inc()
		} else {
			slog.Error("unauthorized", "remoteAddr", r.RemoteAddr, "path", r.URL.Path)
			UnauthorizedTotal.WithLabelValues("", remoteIP, r.URL.Path).Inc()
		}
		errMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		http.Error(w, errMsg, http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, WhoisKey, whois)
	ctx = context.WithValue(ctx, RemoteIPKey, remoteIP)

	m.handler.ServeHTTP(w, r.WithContext(ctx))
}

func CheckRequestPermissionsHandler(localClient *tailscale.LocalClient, accessType string, handler http.Handler) http.Handler {
	return &ACLMiddleware{
		handler:     handler,
		localClient: localClient,
		accessType:  accessType,
	}
}
