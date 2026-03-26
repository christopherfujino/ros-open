package main

import (
	"log"

	"christopherfujino.com/ros/ros-open/notes"
	"christopherfujino.com/ros/ros-open/service"

	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type config struct {
	fs   string
	port int
}

func parseArgs() config {
	var fs = flag.String("fs", "", "Path to mutable file store.")
	var port = flag.Int("port", 8888, "Port.")

	flag.Parse()

	if *fs == "" {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	absFs, err := filepath.Abs(*fs)
	if err != nil {
		panic(err)
	}

	return config{
		fs:   absFs,
		port: *port,
	}
}

func main() {
	var c = parseArgs()

	var services = (func() []service.T {
		var paths = []string{
			"/notes",
		}
		var services = []service.T{}

		for _, path := range paths {
			services = append(services, notes.Create(
				filepath.Join(c.fs, path),
				path,
			))
		}

		return services
	})()

	var descriptions = []service.Description{}
	for _, service := range services {
		service.Register()
		descriptions = append(descriptions, service.Describe())
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<ul>"))
		for _, description := range descriptions {
			fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", description.Endpoint, description.Text)
		}
		w.Write([]byte("</ul>"))
	})

	var localAddress = fmt.Sprintf("127.0.0.1:%d", c.port)

	fmt.Printf("Listening on %s\n", localAddress)
	var err = http.ListenAndServe(
		localAddress,
		loggingMiddleware(http.DefaultServeMux),
	)
	if err != nil {
		panic(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-> %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("<- %s %s", r.Method, r.RequestURI)
	})
}
