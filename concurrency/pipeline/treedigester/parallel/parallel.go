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
package treedigesterparallel

import (
	"crypto/md5"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

// sumFiles starts goroutines to walk the directory tree at root and digest each
// regular file.  These goroutines send the results of the digests on the result
// channel and send the result of the walk on the error channel.  If done is
// closed, sumFiles abandons its work.
func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	// For each regular file, start a goroutine that sums the file and sends
	// the result on c.  Send the result of the walk on errc.
	c := make(chan result)
	errc := make(chan error, 1)

	go func() { // HL
		var wg sync.WaitGroup

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			wg.Add(1)

			go func() { // HL
				data, err := os.ReadFile(path)

				select {
				case c <- result{path, md5.Sum(data), err}: // HL
				case <-done: // HL
				}

				wg.Done()
			}()

			// Abort the walk if done is closed.
			select {
			case <-done: // HL
				return errors.New("walk canceled")
			default:
				return nil
			}
		})

		// Walk has returned, so all calls to wg.Add are done. Start a
		// goroutine to close c once all the sends are done.
		go func() { // HL
			wg.Wait()
			close(c) // HL
		}()

		// No select needed here, since errc is buffered.
		errc <- err // HL
	}()

	return c, errc
}

// md5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, md5All returns an error.  In that case,
// md5All does not wait for inflight read operations to complete.
func md5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{}) // HLdone
	defer close(done)           // HLdone

	resCh, errc := sumFiles(done, root) // HLdone

	m := make(map[string][md5.Size]byte)

	for r := range resCh { // HLrange
		if r.err != nil {
			return nil, r.err
		}

		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
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
