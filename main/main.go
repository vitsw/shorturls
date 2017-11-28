package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"vitsw/urlshort"
)

func main() {
	ymlFile := flag.String("yaml", "", "yaml file with path to url map")
	jsonFile := flag.String("json", "", "json file with path to url map")

	flag.Parse()

	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	mapHandler := urlshort.MapHandler(map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}, mux)

	handler := mapHandler

	if *ymlFile != "" {
		yaml, err := ioutil.ReadFile(*ymlFile);
		if err != nil {
			log.Fatal(err)
		}
		handler, err = urlshort.YAMLHandler(yaml, mapHandler)
		if err != nil {
			panic(err)
		}
	} else if (*jsonFile != "") {
		json, err := ioutil.ReadFile(*jsonFile);
		if err != nil {
			log.Fatal(err)
		}
		handler, err = urlshort.JSONHandler(json, mapHandler)
		if err != nil {
			panic(err)
		}
	} 
	
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
