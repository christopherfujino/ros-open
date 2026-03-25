package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Repo struct {
	Remote string `json:"remote"`
	Description string `json:"description"`
	Local string `json:"local"`
}

type Config struct {
	Repos []Repo `json:"repos"`
}

func Exec(cmdString []string, workingDir string) {
	fmt.Printf("> %s (%s)\n", strings.Join(cmdString, " "), workingDir)
	var cmd = exec.Command(cmdString[0], cmdString[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workingDir
	var err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

// TODO lift to CLI args
var (
	config string = "./sample-config.json"
	repos string = "../repos/remote"
)

func main() {
	flag.StringVar(&repos, "repos", "", "Path to repos directory")
	flag.Parse()

	if repos == "" {
		flag.Usage()
		os.Exit(1)
	}

	bytes, err := os.ReadFile(config)
	if err != nil {
		panic(err)
	}

	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
	reposRoot, err := filepath.Abs(repos)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(reposRoot); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(reposRoot, 0755); err != nil {
			panic(err)
		}
	}

	for _, repo := range config.Repos {
		path, err := filepath.Abs(filepath.Join(reposRoot, repo.Local))
		if err != nil {
			panic(err)
		}
		_, err = os.Stat(path)
		if err != nil {
			Exec([]string{
				"git",
				"clone",
				"--no-single-branch",
				"--tags",
				//"--recurse-submodules",
				"--",
				repo.Remote,
				path,
			}, reposRoot)
			file, err := os.Create(filepath.Join(path, ".git", "description"))
			if err != nil {
				panic(err)
			}
			_, err = file.Write([]byte(repo.Description))
			if err != nil {
				panic(err)
			}
		} else {
			Exec([]string{"git", "pull", "--tags", "--ff-only"}, path)
		}
	}
}
