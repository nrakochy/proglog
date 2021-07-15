package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func newServer() *httpServer {
	return &httpServer {
		Log: NewLog(),
	}
}

func HTTPServer(addr string) *http.Server {
	server := newServer()
	r := mux.NewRouter();
	r.HandleFunc("/", server.handleProduce).Methods("POST")
	r.HandleFunc("/", server.handleConsume).Methods("GET")
	return &http.Server {
		Addr: addr, 
		Handler: r,
	}
}

func (server *httpServer) handleProduce(writer http.ResponseWriter, request *http.Request){
	var req ProduceRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	offset, appendErr := server.Log.Append(req.Record);

	if appendErr != nil {
		http.Error(writer, appendErr.Error(), http.StatusInternalServerError)
		return
	}

	response := ProduceResponse{Offset: offset}
	respErr := json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, respErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *httpServer) handleConsume(writer http.ResponseWriter, request *http.Request){
	var req ConsumeRequest 
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	record, readErr := server.Log.Read(req.Offset);

	if readErr == ErrOffsetNotFound {
		http.Error(writer, readErr.Error(), http.StatusNotFound)
		return
	}

	if readErr != nil {
		http.Error(writer, readErr.Error(), http.StatusInternalServerError)
		return
	}

	response := ConsumeResponse{Record: record}
	respErr := json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, respErr.Error(), http.StatusInternalServerError)
		return
	}
}