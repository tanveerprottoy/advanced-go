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
package treedigesterbounded

import (
	"crypto/md5"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// The MD5All implementation in parallel.go starts a new goroutine for each file.
// In a directory with many large files, this may allocate more memory than is
// available on the machine. We can limit these allocations by bounding the number
// of files read in parallel. In bounded.go, we do this by creating a fixed number
// of goroutines for reading files. Our pipeline now has three stages: walk the tree,
// read and digest the files, and collect the digests.

// The first stage, walkFiles, emits the paths of regular files in the tree:
// walkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel. If done is closed, walkFiles abandons its work.
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() { // HL
		// Close the paths channel after Walk returns.
		defer close(paths) // HL

		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			select {
			case paths <- path: // HL
			case <-done: // HL
				return errors.New("walk canceled")
			}

			return nil
		})
	}()

	return paths, errc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// The middle stage starts a fixed number of digester goroutines that receive file names
// from paths and send results on channel c:
// digester reads path names from paths and sends digests of the corresponding
// files on c until either paths or done is closed.
func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
	for path := range paths { // HLpaths
		data, err := os.ReadFile(path)

		select {
		case c <- result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

// md5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, md5All returns an error.  In that case,
// md5All does not wait for inflight read operations to complete.
func md5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan result) // HLc

	var wg sync.WaitGroup
	
	// this limit the number of goroutines to 20
	const numDigesters = 20

	wg.Add(numDigesters)

	for range numDigesters {
		go func() {
			digester(done, paths, c) // HLc

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()

		close(c) // HLc
	}()
	// End of pipeline

	m := make(map[string][md5.Size]byte)

	// The final stage receives all the results from c then checks the error from errc.
	// This check cannot happen any earlier, since before this point, walkFiles may block
	// sending values downstream:
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}

		m[r.path] = r.sum
	}

	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		return nil, err
	}

	return m, nil
}

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
