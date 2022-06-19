package raft

import (
	"context"
	"github.com/roessland/raft-consensus/raft/sets"
	"github.com/roessland/raft-consensus/raft/stable"
	"math/rand"
	"net/http"
	"time"
)

const NumNodes = 3

type LogEntry struct {
	// The term the entry was stored in
	term int

	// The log data sent by the client
	command []byte

	// The index the entry is stored at
	index int
}

type Node struct {
	nodeId   int // identity of the current node.
	numNodes int // total number of raft nodes.

	currentTerm   stable.IntStore
	votedFor      stable.NullableIntStore
	log           stable.LogEntriesStore
	commitLength  stable.IntStore
	currentRole   Role
	currentLeader *int
	votesReceived sets.IntSet
	sentLength    [NumNodes]int
	ackedLength   [NumNodes]int

	// Selectable conditions
	incomingVoteRequests      chan VoteRequestParams
	electionTimeoutEvents     chan struct{} // For candidates
	leaderTimeoutEvents       chan struct{} // For followers
	shouldSendHeartbeatEvents chan struct{} // For leader

	// sendHeartbeatTimer TODO
	// suspectsLeaderFailureTimer TODO
	// replicateLogTimer TODO https://youtu.be/uXEYuDwm7e4?t=952

	// Timers
	electionTimer              *time.Timer   // For candidates
	heartbeatTimer             *time.Timer   // For followers
	shouldSendHeartbeatTimer   *time.Timer   // For leader
	electionTimeout            time.Duration // For candidates
	heartbeatTimeout           time.Duration // For followers
	shouldSendHeartbeatTimeout time.Duration // For leader

	// Non-raft stuff
	httpClient *http.Client
	done       chan struct{}
}

func NewNode(nodeId int) *Node {
	n := &Node{}
	n.numNodes = 3
	n.nodeId = nodeId

	if !n.needsInitialisation() {
		n.initStableVariables()
	}

	n.initTransientVariables()

	n.initTimers()

	n.httpClient = &http.Client{}
	n.done = make(chan struct{})

	return n
}

// needsInitialisation checks if this is the first run,
// and stable variables must be initialized, or if this is
// a crash, where stable variables already exist on disk.
func (n *Node) needsInitialisation() bool {
	// depends on commitLength being set last in initStableVariables.
	return !n.commitLength.AlreadyExists()
}

func (n *Node) initStableVariables() {
	n.currentTerm.Set(0)
	n.votedFor.SetNull()
	n.log.SetEmpty()
	n.commitLength.Set(0)
}

func (n *Node) initTransientVariables() {
	n.currentRole = roleFollower
	n.currentLeader = nil
	n.votesReceived = sets.NewIntSet()
	n.sentLength = [3]int{0, 0, 0}
	n.ackedLength = [3]int{0, 0, 0}
}

func (n *Node) initTimers() {
	n.electionTimer = time.NewTimer(randomInterval())
	n.heartbeatTimer = time.NewTimer(randomInterval())
	n.shouldSendHeartbeatTimer = time.NewTimer(randomInterval())
}

func (n *Node) Start() {
	go n.serveRPC()
}

func (n *Node) Close() {
	close(n.done)
}

func (n *Node) Command(ctx context.Context, cmd []byte) error {
	return nil
}

func randomInterval() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}
