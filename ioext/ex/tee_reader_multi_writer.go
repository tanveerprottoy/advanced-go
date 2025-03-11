package ex

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func teeReaderMultiWriter(w http.ResponseWriter, r *http.Request) {
	req := make(map[string]any)

	// tee reader
	json.NewDecoder(io.TeeReader(r.Body, os.Stdout)).
		Decode(&req)

	response := map[string]string{
		"message": "Looking great",
	}

	// multi writer
	json.NewEncoder(io.MultiWriter(os.Stdout, w)).
		Encode(response)
}

func ExecuterTeeReaderMultiWriter() {
	http.HandleFunc("/", teeReaderMultiWriter)
	http.ListenAndServe(":8080", nil)
}
