# kimchi

A bare-bones HTTP server. Designed to be used together with [tlstunnel].

```
site example.org {
	root /srv/http
}

site example.com {
	proxy http://localhost:8080
}
```

## License

MIT

[tlstunnel]: https://sr.ht/~emersion/tlstunnel
