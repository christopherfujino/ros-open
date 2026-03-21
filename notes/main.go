package notes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"christopherfujino.com/ros/ros-open/service"
)

type GetNotesResponse struct {
	Files []string `json:"files"`
}

type GetNoteResponse struct {
	Content string `json:"content"`
}

type UpdateNotesRequest struct {
	Path     string `json:"path"`
	Contents string `json:"contents"`
}

type tee struct {
	dir          string
	endpointRoot string
}

func (t tee) Describe() service.Description {
	return service.Description{
		Endpoint: t.endpointRoot,
		Text: "Notes",
	}
}

func (t tee) Register() {
	// TODO fix this
	assetsPath, err := filepath.Abs(filepath.Join(".", "notes", "assets"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("note asset path = %s\n", assetsPath)
	var fileServer = http.StripPrefix(t.endpointRoot, http.FileServer(http.Dir(assetsPath)))
	http.Handle(fmt.Sprintf("GET %s/", t.endpointRoot), fileServer)

	// GET one
	http.HandleFunc(fmt.Sprintf("GET /api%s/note/{name...}", t.endpointRoot), func(w http.ResponseWriter, r *http.Request) {
		var fail = func(err error) {
			w.WriteHeader(500)
			// TODO sanitize this?
			_, _ = w.Write([]byte(err.Error()))
			log.Printf("Error: %v", err)
		}

		// Add back leading backslash
		var localName = r.PathValue("name")
		if localName == "" {
			panic(fmt.Sprintf("Huh?! %s\n", r.URL))
		}
		fs, err := Open(t.dir)
		if err != nil {
			fail(err)
			return
		}
		bytes, err := fs.ReadFile(localName)
		if err != nil {
			fail(err)
			return
		}
		var res = GetNoteResponse{
			Content: string(bytes),
		}
		bytes, err = json.Marshal(res)
		if err != nil {
			fail(err)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			log.Panicf("Failed to write response: %s", err.Error())
		}
	})
	// GET all
	http.HandleFunc(fmt.Sprintf("GET /api%s/notes", t.endpointRoot), func(w http.ResponseWriter, r *http.Request) {
		var fail = func(err error) {
			w.WriteHeader(500)
			// TODO sanitize this?
			_, _ = w.Write([]byte(err.Error()))
		}
		var fs, err = Open(t.dir)
		if err != nil {
			fail(err)
			return
		}
		log.Printf("GET %v\n", r.URL.Path)

		files, err := fs.GetAllPaths()
		if err != nil {
			fail(err)
			return
		}
		var res = GetNotesResponse{
			Files: files,
		}
		resBytes, err := json.Marshal(res)
		if err != nil {
			fail(err)
			return
		}
		_, err = w.Write(resBytes)
		if err != nil {
			log.Printf("Error failed to write response: %s\n", err.Error())
		}
	})

	http.HandleFunc("UPDATE /api/notes/update", func(w http.ResponseWriter, r *http.Request) {
		var fail = func(err error) {
			w.WriteHeader(500)
			// TODO sanitize this?
			_, _ = w.Write([]byte(err.Error()))
		}

		var fs, err = Open(t.dir)
		if err != nil {
			fail(err)
			return
		}
		log.Printf("UPDATE %v\n", r.URL.Path)
		var buffer = bytes.Buffer{}
		io.Copy(&buffer, r.Body)

		var reqData UpdateNotesRequest
		err = json.Unmarshal(buffer.Bytes(), &reqData)
		if err != nil {
			w.WriteHeader(500)
			// TODO sanitize this?
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		fs.Write(reqData.Path, reqData.Contents)
		w.WriteHeader(200)
	})
}

func Create(dirPath, endpointRoot string) service.T {
	return tee{
		dir: dirPath,
		endpointRoot: endpointRoot,
	}
}
