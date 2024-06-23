package progress

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type progress struct {
	total uint64
}

func (p *progress) Write(b []byte) (int, error) {
	p.total += uint64(len(b))
	fmt.Printf("Downloaded %d bytes...\n", p.total)
	return len(b), nil
}

func Executer() {
	res, err := http.Get("http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-5gram-20120701-0.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	localFile, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer localFile.Close()

	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	teeReader := io.TeeReader(gzipReader, &progress{})

	if _, err := io.Copy(localFile, teeReader); err != nil {
		log.Fatal(err)
	}
}
