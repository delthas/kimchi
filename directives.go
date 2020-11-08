package main

import (
	"fmt"

	"git.sr.ht/~emersion/go-scfg"
)

func parseConfig(srv *Server, cfg scfg.Block) error {
	for _, dir := range cfg {
		switch dir.Name {
		case "site":
			if err := parseSite(srv, dir); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown directive %q", dir.Name)
		}
	}
	return nil
}

func parseSite(srv *Server, dir *scfg.Directive) error {
	for _, addr := range dir.Params {
		srv.AddListener("tcp", addr)
	}
	return nil
}
