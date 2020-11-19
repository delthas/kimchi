package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

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
	for _, uriStr := range dir.Params {
		if !strings.Contains(uriStr, "//") {
			uriStr = "//" + uriStr
		}

		u, err := url.Parse(uriStr)
		if err != nil {
			return err
		}

		var ln *Listener
		var host, port string
		switch u.Scheme {
		case "", "http", "http+insecure":
			if host, port, err = net.SplitHostPort(u.Host); err != nil {
				host = u.Host
				port = "http"
			}
			ln = srv.AddListener("tcp", ":"+port)
			if u.Scheme == "http+insecure" {
				ln.Insecure = true
			}
		default:
			return fmt.Errorf("unknown URI scheme %q", u.Scheme)
		}

		path := u.Path
		if path == "" {
			path = "/"
		}
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("invalid path %q", path)
		}

		pattern := host+path

		// First process handler directives
		var handler http.Handler
		for _, child := range dir.Children {
			var h http.Handler
			switch child.Name {
			case "root":
				var dir string
				if err := child.ParseParams(&dir); err != nil {
					return err
				}
				h = http.FileServer(http.Dir(dir))
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
				h = proxy
			}
			if h != nil {
				if handler != nil {
					return fmt.Errorf("multiple HTTP handler directives provided")
				}
				handler = h
			}
		}
		if handler == nil {
			return fmt.Errorf("missing handler directive")
		}

		// Then process middleware directives
		for _, child := range dir.Children {
			switch child.Name {
			case "root", "reverse_proxy":
				// Handler directive already processed above
			default:
				handler, err = parseMiddleware(child, handler)
				if err != nil {
					return fmt.Errorf("directive %q: %v", child.Name, err)
				}
			}
		}

		ln.Mux.Handle(pattern, handler)
	}
	return nil
}

func parseMiddleware(dir *scfg.Directive, next http.Handler) (http.Handler, error) {
	switch dir.Name {
	case "header":
		// TODO: allow adding and removing fields
		setFields := make(map[string]string)
		if len(dir.Params) > 0 {
			if len(dir.Params) != 2 {
				return nil, fmt.Errorf("expected exactly two parameters")
			}
			setFields[dir.Params[0]] = dir.Params[1]
		} else {
			for _, child := range dir.Children {
				if len(child.Params) != 1 {
					return nil, fmt.Errorf("expected exactly one parameter for child directive")
				}
				if _, ok := setFields[child.Name]; ok {
					return nil, fmt.Errorf("duplicate child directive %q", child.Name)
				}
				setFields[child.Name] = child.Params[0]
			}
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range setFields {
				w.Header().Set(k, v)
			}
			next.ServeHTTP(w, r)
		}), nil
	default:
		return nil, fmt.Errorf("unknown directive")
	}
}
