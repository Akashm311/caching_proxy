package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"io"

)

type CachedResponse struct {
    Body    []byte
    Headers http.Header
}

func main() {
	// Define command-line flags
	port := flag.Int("port", 8080, "port to run the server on")
	origin := flag.String("origin", "http://localhost:3000", "allowed origin")
	clear := flag.Bool("clear-cache", false, "clear the cache")
	cache := make(map[string]CachedResponse)
	

	// Parse the command-line flags
	flag.Parse()

	if *clear {
		// Logic to clear the cache
		fmt.Println("Cache cleared.")
		return
	} else {
		// validate that port and origin are provided
		if *port <= 0 || *port > 65535 {
			log.Fatalf("Invalid port number: %d. Port must be between 1 and 65535.", *port)
		}
		if *origin == "" {
			log.Fatal("Origin cannot be empty.")
		}
	}

	// Use the parsed flags
	fmt.Printf("Server will run on port %d.\n", *port)
	fmt.Printf("Allowed origin: %s.\n", *origin)

	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		targetURL := fmt.Sprintf("%s%s", *origin, r.URL.Path)

		if r.URL.Path == "/clear-cache" && r.Method == "DELETE" {
			cache = make(map[string]CachedResponse)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Cache cleared"))
			return
		}

		if cachedResponse, found := cache[r.URL.Path]; found {
			log.Printf("Cache hit for %s", r.URL.Path)
			for key, values := range cachedResponse.Headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Set("X-Cache", "HIT")
			w.Write(cachedResponse.Body)
			return
		}

		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, "Failed to reach origin server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		cache[r.URL.Path] = CachedResponse{
			Body:    body,
			Headers: resp.Header,
		}
		w.Header().Set("X-Cache", "MISS")
		w.Write(body)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}