package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	// Define command-line flags
	dirPath := flag.String("dir", ".", "directory path to serve")
	username := flag.String("username", "guest", "basic auth username")
	password := flag.String("password", "guest", "basic auth password")
	port := flag.Int("port", 8888, "server port number")
	flag.Parse()

	// Create a file server handler for the specified directory
	fs := http.FileServer(http.Dir(*dirPath))

	// Create a middleware that checks the Authorization header for valid credentials
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if *username != "" && *password != "" {
				u, p, ok := r.BasicAuth()
				if !ok || u != *username || p != *password {
					w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintln(w, "Unauthorized")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}

	// Wrap the file server handler with the auth middleware
	http.Handle("/", authMiddleware(fs))

	// Start the server
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server listening on %s, serving directory %q\n", addr, *dirPath)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
