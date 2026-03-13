package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/strugglerx/sjson"
)

type expandRequest struct {
	Data json.RawMessage `json:"data"`
}

type expandResponse struct {
	Result string      `json:"result,omitempty"`
	Parsed interface{} `json:"parsed,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/expand", expandHandler)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, expandResponse{
			Error: "method not allowed",
		})
		return
	}

	defer r.Body.Close()

	var req expandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, expandResponse{
			Error: "invalid request body",
		})
		return
	}

	result, err := sjson.StringWithJsonScanToStringE(req.Data)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, expandResponse{
			Error: err.Error(),
		})
		return
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		writeJSON(w, http.StatusInternalServerError, expandResponse{
			Error: "expanded result is not valid json",
		})
		return
	}

	writeJSON(w, http.StatusOK, expandResponse{
		Result: result,
		Parsed: parsed,
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}
