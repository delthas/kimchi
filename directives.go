package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

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
			case "reverse_proxy":
				var urlStr string
				if err := child.ParseParams(&urlStr); err != nil {
					return err
				}
				target, err := url.Parse(urlStr)
				if err != nil {
					return err
				}
				proxy := httputil.NewSingleHostReverseProxy(target)
				director := proxy.Director
				proxy.Director = func(req *http.Request) {
					director(req)
					req.Host = target.Host
				}
				ln.Mux.Handle("/", proxy)
			default:
				return fmt.Errorf("unknown directive %q", child.Name)
			}
		}
	}
	return nil
}
