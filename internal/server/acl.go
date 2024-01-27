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
	Read  []string `json:"read"`
	Write []string `json:"write"`
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
	whois, err := localClient.WhoIs(r.Context(), r.RemoteAddr)
	if err != nil {
		return whois, "", err
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return whois, host, err
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

func (m *ACLMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	whois, remoteIP, err := CheckRequestPermissions(m.localClient, r, m.accessType)
	if err != nil {
		if whois != nil {
			slog.Error("unauthorized", "remoteAddr", r.RemoteAddr, "path", r.URL.Path)
			UnauthorizedTotal.WithLabelValues("", "", r.URL.Path).Inc()
		} else {
			slog.Warn("unauthorized", "hostname", whois.Node.Name, "ip", remoteIP, "path", r.URL.Path)
			UnauthorizedTotal.WithLabelValues(whois.Node.Name, remoteIP, r.URL.Path).Inc()
		}
		errMsg := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		http.Error(w, errMsg, http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), WhoisKey, whois)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}

func CheckRequestPermissionsHandler(localClient *tailscale.LocalClient, accessType string, handler http.Handler) http.Handler {
	return &ACLMiddleware{
		handler:     handler,
		localClient: localClient,
		accessType:  accessType,
	}
}
