package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roessland/raft-consensus/raft"
	"log"
	"net/http"
)

var nodeId int
var raftNode *raft.Node

func init() {
	flag.IntVar(&nodeId, "nodeid", 0, "must be unique per process")
}

func main() {
	flag.Parse()

	raftNode = raft.NewNode(nodeId)

	r := mux.NewRouter()
	r.HandleFunc("/set/{key}/{val}", handleSet)
	r.HandleFunc("/get/{key}", handleGet)

	addr := fmt.Sprintf("127.0.0.1:%d", 8000+nodeId)
	log.Printf("API: Listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["val"]
	if key == "" || val == "" {
		http.Error(w, "missing key or value", http.StatusBadRequest)
		return
	}

	err := raftNode.Broadcast(r.Context(), []byte(fmt.Sprintf(`SET %s %s`, key, val)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {

}
