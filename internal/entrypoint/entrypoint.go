package entrypoint

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"time"

	"github.com/mastrasec/firewow/internal/service"
)

type Entry struct {
	svc    service.Contract
	server *http.Server
	mux    *http.ServeMux
}

func New(addr string, svc service.Contract) *Entry {
	mux := http.NewServeMux()

	return &Entry{
		svc: svc,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			WriteTimeout: 10 * time.Second,
		},
		mux: mux,
	}
}

func (entry *Entry) oaiV1ChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	header := map[string][]string(r.Header)
	body, err := getBody(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	resHeaders, resBody, err := entry.svc.OpenAIV1ChatCompletionsPipe(header, body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for key, values := range resHeaders {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	res, err := gzipCompression(resBody)
	if err != nil {
		return
	}

	if _, err := w.Write(res); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (entry *Entry) RegisterRoutes() {
	entry.mux.HandleFunc("/oai/v1/chat/completions", entry.oaiV1ChatCompletions)
}

func (entry *Entry) Run() error {
	entry.RegisterRoutes()

	return entry.server.ListenAndServe()
}

func getBody(data io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	return body, nil
}

func gzipCompression(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
