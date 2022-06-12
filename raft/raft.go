package raft

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const NumNodes = 3

type Role int

const (
	follower  Role = iota
	leader    Role = iota
	candidate Role = iota
)

type logEntry struct {
	// The term the entry was stored in
	currentTerm int

	// The log data sent by the client
	command []byte

	// The index the entry is stored at
	index int
}

type VoteRequest struct {
	CandidateID        int
	CandidateTerm      int
	CandidateLogLength int
	CandidateLogTerm   int
}

type Node struct {
	nodeId               int
	currentTerm          int
	votedFor             int
	log                  []logEntry
	commitLength         int
	currentRole          Role
	currentLeader        int
	votesReceived        map[int]struct{}
	sentLength           [NumNodes]int
	ackedLength          [NumNodes]int
	incomingVoteRequests chan VoteRequest
}

func NewNode(nodeId int) *Node {
	n := &Node{}
	n.votedFor = -1
	n.currentLeader = -1
	n.nodeId = nodeId

	go n.serveRPC()
	return n
}

func (n *Node) Command(ctx context.Context, cmd []byte) error {
	return nil
}

func (n *Node) httpMsgHandler(w http.ResponseWriter, r *http.Request) {

}

func (n *Node) serveRPC() {
	r := mux.NewRouter()
	r.HandleFunc("/", n.httpMsgHandler)
	addr := fmt.Sprintf("127.0.0.1:%d", 50000+n.nodeId)
	log.Printf("Raft: listening on HTTP at %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
