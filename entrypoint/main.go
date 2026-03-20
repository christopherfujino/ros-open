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

func parseArgs() string {
	var fs = flag.String("fs", "", "Path to mutable file store.")

	flag.Parse()

	if *fs == "" {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	absFs, err := filepath.Abs(*fs)
	if err != nil {
		panic(err)
	}

	*fs = absFs

	return *fs
}

func main() {
	var storageRoot = parseArgs()

	var services = (func () []service.T {
		var paths = []string{
			"/notes",
		}
		var services = []service.T{}

		for _, path := range paths {
			services = append(services, notes.Create(
				filepath.Join(storageRoot, path),
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
	fmt.Println("Listening on 127.0.0.1:8080")

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Debug in GET /: %s\n", r.URL.String())
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<ul>"))
		for _, description := range descriptions {
			fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", description.Endpoint, description.Text)
		}
		w.Write([]byte("</ul>"))
	})

	var err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}
