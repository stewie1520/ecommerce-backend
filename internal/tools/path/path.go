package path

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

func PathCWD(pathname string) string {
	return path.Join(CWD(), pathname)
}


func CWD() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panicln(err)
	}

	dataDir := path.Join(dir, "..")
	return dataDir
}
