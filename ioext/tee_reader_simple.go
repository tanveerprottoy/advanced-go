package ioext

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
)

/*
tee reader
  ... | step 0 | step 1 | step 3 | ...

>----->------>----->------┬------>------>------>----->
[stdin]                   |                   [stdout]
                          |
                          |
                          v
                       [log.log]

*/

func hashAndSendSimple(r io.Reader) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		log.Println(fmt.Errorf("%v", err))
	}
	log.Println(string(bytes))

	h := sha256.New()
	h.Write(bytes)
	sha := hex.EncodeToString(h.Sum(nil))
	log.Println(sha)
}

func hashAndSend(r io.Reader) {
	w := sha256.New()

	// any reads from tee will read from r and write to w
	tee := io.TeeReader(r, w)

	sendReader(tee)

	sha := hex.EncodeToString(w.Sum(nil))
	fmt.Println(sha)
}

// sendReader sends the contents of an io.Reader to stdout using a 256 byte buffer
func sendReader(data io.Reader) {
	buff := make([]byte, 256)
	for {
		_, err := data.Read(buff)
		if err == io.EOF {
			break
		}

		log.Print(string(buff))
	}

	log.Println("")
}

func Executer() {
	r1 := strings.NewReader("hello world")
	r2 := strings.NewReader("hello world")

	hashAndSendSimple(r1)
	hashAndSend(r2)
}
