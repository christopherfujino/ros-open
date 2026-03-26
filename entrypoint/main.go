package main

import (
	"log"

	"christopherfujino.com/ros/ros-open/globals"
	"christopherfujino.com/ros/ros-open/notes"
	"christopherfujino.com/ros/ros-open/service"

	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func parseArgs() globals.T {
	var fs = flag.String("fs", "", "Path to mutable file store.")
	var repoRoot = flag.String("src", "", "Path to the root of the ros-open repo.")
	var port = flag.Int("port", 8888, "Port.")

	flag.Parse()

	var failed = false

	if *fs == "" {
		flag.CommandLine.Usage()
		failed = true
	}

	if *repoRoot == "" {
		flag.CommandLine.Usage()
		failed = true
	}

	if failed {
		os.Exit(1)
	}

	absFs, err := filepath.Abs(*fs)
	if err != nil {
		panic(err)
	}

	return globals.T{
		FileStoreRoot: absFs,
		Port:          *port,
		RosOpenRoot:   *repoRoot,
	}
}

func main() {
	var c = parseArgs()

	type tuple struct {
		endpointPath string
		registrar    func(globals.T, string) service.T
	}

	var services = (func() []service.T {
		var paths = []tuple{
			tuple{"/notes", notes.Create},
		}
		var services = []service.T{}

		for _, t := range paths {
			services = append(services, t.registrar(
				c,
				t.endpointPath,
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

	var localAddress = fmt.Sprintf("127.0.0.1:%d", c.Port)

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
