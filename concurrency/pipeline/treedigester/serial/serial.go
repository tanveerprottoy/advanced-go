/*
MD5 is a message-digest algorithm that’s useful as a file checksum.
The command line utility md5sum prints digest values for a list of files.

% md5sum *.go
d47c2bbc28298ca9befdfbc5d3aa4e65  bounded.go
ee869afd31f83cbb2d10ee81b2b831dc  parallel.go
b88175e65fdcbc01ac08aaf1fd9b5e96  serial.go

Our example program is like md5sum but instead takes a single
directory as an argument and prints the digest values for each
regular file under that directory, sorted by path name.

% go run serial.go .
d47c2bbc28298ca9befdfbc5d3aa4e65  bounded.go
ee869afd31f83cbb2d10ee81b2b831dc  parallel.go
b88175e65fdcbc01ac08aaf1fd9b5e96  serial.go
*/

package treedigesterserial

import (
	"crypto/md5"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// md5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.
// The MD5All function is the focus of our discussion. this implementation
// uses no concurrency and simply reads and sums each file as it walks the tree.
func md5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		data, err := os.ReadFile(path) // HL
		if err != nil {
			return err
		}

		m[path] = md5.Sum(data) // HL
		
		return nil
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

// The Executer function of our program invokes a helper function
// md5All, which returns a map from path name to digest value,
// then sorts and prints the results
func Executer() {
	// Calculate the MD5 sum of all files under the specified directory,
	// then print the results sorted by path name.
	m, err := md5All(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	for _, path := range paths {
		log.Printf("%x  %s\n", m[path], path)
	}
}
