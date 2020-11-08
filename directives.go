package main

import (
	"git.sr.ht/~emersion/go-scfg"
)

func parseSite(srv *Server, dir *scfg.Directive) error {
	for _, addr := range dir.Params {
		srv.AddListener("tcp", addr)
	}
	return nil
}
