package main

import (
	"christopherfujino.com/ros/ros-open/notes"
	"christopherfujino.com/ros/ros-open/service"

	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type config struct {
	fs string
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
		fs: absFs,
		port: *port,
	}
}

func main() {
	var c = parseArgs()

	var services = (func () []service.T {
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
		fmt.Printf("Debug in GET /: %s\n", r.URL.String())
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<ul>"))
		for _, description := range descriptions {
			fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", description.Endpoint, description.Text)
		}
		w.Write([]byte("</ul>"))
	})

	fmt.Printf("Listening on 0.0.0.0:%d\n", c.port)
	var err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", c.port), nil)
	if err != nil {
		panic(err)
	}
}
