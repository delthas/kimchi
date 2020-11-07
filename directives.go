package main

import (
	"git.sr.ht/~emersion/go-scfg"
)

func parseSite(srv *Server, dir *scfg.Directive) error {
	for _, addr := range dir.Params {
		srv.Listeners = append(srv.Listeners, &Listener{
			Network: "tcp",
			Address: addr,
			Server:  srv,
		})
	}
	return nil
}
