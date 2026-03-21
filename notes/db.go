package notes

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// TODO use either sqlite or just the file system
// TODO rename DB -> FS

type FS struct {
	root string
}

func Open(root string) (*FS, error) {
	var abs, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	// TODO check it's a valid directory
	return &FS{
		root: abs,
	}, nil
}

func (f FS) ReadFile(relativePath string) ([]byte, error) {
	var bytes, err = os.ReadFile(filepath.Join(f.root, relativePath))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (f *FS) Write(path string, contents string) {
	log.Printf("FS.Write(%s, ...)", path)
	var absPath = filepath.Join(f.root, path)
	_, err := os.Stat(absPath)
	if err == nil {
		fmt.Fprintf(os.Stderr, "About to overwrite %s\n", path)
	} else if errors.Is(err, os.ErrNotExist) {
		// TODO we must create parents

		var findFirstExistingParent func(string) = nil
		findFirstExistingParent = func(path string) {
			if path == "/" {
				panic("Uh oh, traversed too far up!")
			}
			var parent = filepath.Dir(path)
			_, err := os.Stat(parent)
			if err == nil {
				// We found it!
				return
			} else if errors.Is(err, os.ErrNotExist) {
				findFirstExistingParent(parent)
				err = os.Mkdir(parent, 0700)
				if err != nil {
					panic(err)
				}
				log.Printf("Created directory %s\n", parent)
			}
		}
		findFirstExistingParent(absPath)
	} else {
		panic(err)
	}
	if err = os.WriteFile(absPath, []byte(contents), 0600); err != nil {
		panic(err)
	}
	log.Printf("Wrote file %s\n", absPath)
}

func (f FS) GetAllPaths() ([]string, error) {
	var paths = []string{}

	var iterateDir func(string, string) error
	iterateDir = func(absoluteParent string, relativeParent string) error {
		var dirEntries, err = os.ReadDir(absoluteParent)
		if err != nil {
			return err
		}

		for _, dirEntry := range dirEntries {
			var absoluteName = filepath.Join(absoluteParent, dirEntry.Name())
			var relativeName = filepath.Join(relativeParent, dirEntry.Name())
			if !dirEntry.IsDir() {
				log.Printf("Found file %s\n", relativeName)
				// TODO we should only append as far as to the f.root
				paths = append(paths, relativeName)
			} else {
				log.Printf("Found dir %s\n", relativeName)
				iterateDir(absoluteName, relativeName)
			}
		}

		return nil
	}

	iterateDir(f.root, "/")

	return paths, nil
}
