package main

import (
	"net/http"
	"fmt"
	"sort"

	"github.com/dunglas/httpsfv"
)

type acceptEncoding struct {
	Name   string
	Weight float64
}

func unmarshalAcceptEncoding(h http.Header) ([]acceptEncoding, error) {
	l, err := httpsfv.UnmarshalList(h["Accept-Encoding"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse Accept-Encoding: %v", err)
	}

	var encodings []acceptEncoding
	for _, member := range l {
		item, ok := member.(httpsfv.Item)
		if !ok {
			return nil, fmt.Errorf("failed to parse Accept-Encoding: expected Item, got %T", member)
		}
		name, ok := item.Value.(string)
		if !ok {
			return nil, fmt.Errorf("failed to parse Accept-Encoding: expected item value to be string, got %T", item.Value)
		}
		var weight float64
		if v, ok := item.Params.Get("q"); ok {
			weight, ok = v.(float64)
			if !ok {
				return nil, fmt.Errorf("failed to parse Accept-Encoding: expected \"q\" value to be float64, got %T", v)
			}
		}
		encodings = append(encodings, acceptEncoding{
			Name:   name,
			Weight: weight,
		})
	}

	sort.Slice(encodings, func(i, j int) bool {
		return encodings[i].Weight > encodings[j].Weight
	})

	return encodings, nil
}
