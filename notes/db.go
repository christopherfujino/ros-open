package notes

// TODO use either sqlite or just the file system

type DB map[string]string

var global DB

func Open() *DB {
	//return &DB{}
	return &global
}

func (d *DB) Write(path string, contents string) {
	var m = map[string]string(*d)
	m[path] = contents
}

func (d DB) GetAllPaths() []string {
	var paths []string

	for k, _ := range d {
		paths = append(paths, k)
	}

	return paths
}
