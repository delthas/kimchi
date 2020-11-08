package main

import (
	"fmt"
	"net/http"

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
		ln := srv.AddListener("tcp", addr)

		for _, child := range dir.Children {
			switch child.Name {
			case "root":
				var dir string
				if err := child.ParseParams(&dir); err != nil {
					return err
				}
				ln.Mux.Handle("/", http.FileServer(http.Dir(dir)))
			default:
				return fmt.Errorf("unknown directive %q", child.Name)
			}
		}
	}
	return nil
}
