package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Repo struct {
	Remote string `json:"remote"`
	Description string `json:"description"`
	Local string `json:"local"`
}

type Config struct {
	Repos []Repo `json:"repos"`
}

func main() {
	bytes, err := os.ReadFile("./sample-config.json")
	if err != nil {
		panic(err)
	}

	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
	const reposRoot = "../repos"
	for _, repo := range config.Repos {
		path, err := filepath.Abs(filepath.Join(reposRoot, repo.Local))
		if err != nil {
			panic(err)
		}
		_, err = os.Stat(path)
		if err != nil {
			panic("TODO git clone...")
		} else {
			fmt.Printf("No error stat-ing %s\n", path)
		}
		fmt.Println(path)
		fmt.Printf("%v\n", repo)
	}
}
