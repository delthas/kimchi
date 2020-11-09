# [kimchi]

A bare-bones HTTP server. Designed to be used together with [tlstunnel].

```
site example.org {
	root /srv/http
}

site example.com {
	reverse_proxy http://localhost:8080
}
```

## License

MIT

[kimchi]: https://sr.ht/~emersion/kimchi
[tlstunnel]: https://sr.ht/~emersion/tlstunnel
