package server

import (
	"errors"
	"net/http"
	"path/filepath"

	"tailscale.com/client/tailscale"
	"tailscale.com/tailcfg"
)

const PeerCapName = "github.com/josh/datajar-server"

type Capabilities struct {
	Read  []string `json:"read"`
	Write []string `json:"write"`
}

func CanAccessPath(path string, caps []Capabilities, accessType string) bool {
	for _, cap := range caps {
		var accessList []string
		if accessType == "read" {
			accessList = cap.Read
		} else if accessType == "write" {
			accessList = cap.Write
		}

		for _, access := range accessList {
			if access == "*" {
				return true
			}

			allowed, _ := filepath.Match("/"+access, path)
			if allowed {
				return true
			}
		}
	}
	return false
}

func CheckRequestPermissions(localClient *tailscale.LocalClient, r *http.Request, accessType string) error {
	whois, err := localClient.WhoIs(r.Context(), r.RemoteAddr)
	if err != nil {
		return err
	}

	caps, err := tailcfg.UnmarshalCapJSON[Capabilities](whois.CapMap, PeerCapName)
	if err != nil {
		return err
	}

	if !CanAccessPath(r.URL.Path, caps, accessType) {
		return errors.New("unauthorized")
	}

	return nil
}
