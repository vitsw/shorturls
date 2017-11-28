package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildPathMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     [{"path": "/some-path", "url": "https://www.some-url.com/demo"}]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildPathMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

type redirects []struct {
	Path string
	URL  string
}

func parseYaml(yml []byte) (redirects, error) {
	var urlsMap redirects
	if err := yaml.Unmarshal(yml, &urlsMap); err != nil {
		return nil, err
	}
	return urlsMap, nil
}

func parseJSON(jsonStr []byte) (redirects, error) {
	var pathToUrls redirects
	if err := json.Unmarshal(jsonStr, &pathToUrls); err != nil {
		return nil, err
	}
	return pathToUrls, nil
}

func buildPathMap(pMap redirects) map[string]string {
	pathMap := make(map[string]string, len(pMap))
	for _, item := range pMap {
		pathMap[item.Path] = item.URL
	}
	return pathMap
}
