package stream

import (
	"encoding/json"
	"net/http"
)

type ResponseWriterFlusher interface {
	http.ResponseWriter
	http.Flusher
}

type StreamWriter interface {
	Write(v any) error
}

type defaultWriter struct {
	res     ResponseWriterFlusher
	encoder *json.Encoder
}

func New(fw ResponseWriterFlusher) StreamWriter {
	fw.Header().Set("Content-Type", "application/json")
	fw.Header().Set("Transfer-Encoding", "chunked")
	fw.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(fw)
	return &defaultWriter{res: fw, encoder: encoder}
}

func (w *defaultWriter) Write(v any) error {
	if err := w.encoder.Encode(v); err != nil {
		return err
	}
	w.res.Flush()
	return nil
}

type message struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

func Error(w StreamWriter, payload any) error {
	return Write(w, "error", payload)
}

func Message(w StreamWriter, payload any) error {
	return Write(w, "message", payload)
}

func Write(w StreamWriter, messageType string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return w.Write(message{Type: messageType, Message: data})
}
