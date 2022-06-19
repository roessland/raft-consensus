package kv

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roessland/raft-consensus/raft"
	"log"
	"net/http"
)

type Server struct {
	nodeId   int
	raftNode *raft.Node
	router   *mux.Router
	addr     string
}

func NewServer(nodeId int) *Server {
	s := &Server{}
	s.nodeId = nodeId
	s.raftNode = raft.NewNode(nodeId)

	r := mux.NewRouter()
	r.HandleFunc("/set/{key}/{val}", s.handleSet)
	r.HandleFunc("/get/{key}", s.handleGet)
	s.router = r

	s.addr = fmt.Sprintf("127.0.0.1:%d", 8000+nodeId)

	return s
}

func (s *Server) ListenAndServe() {
	log.Printf("API: Listening at %s", s.addr)
	log.Print(http.ListenAndServe(s.addr, s.router))
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["val"]
	if key == "" || val == "" {
		http.Error(w, "missing key or value", http.StatusBadRequest)
		return
	}

	err := s.raftNode.Broadcast(r.Context(), []byte(fmt.Sprintf(`SET %s %s`, key, val)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	// ask all other nodes what they think
}
