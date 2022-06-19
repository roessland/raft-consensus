package raft

import (
	"github.com/roessland/raft-consensus/raft/sets"
	"github.com/roessland/raft-consensus/raft/stable"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const NumNodes = 3

type StableStorage struct {
	CurrentTerm  stable.IntStore
	VotedFor     stable.NullableIntStore
	Log          stable.LogEntriesStore
	CommitLength stable.IntStore
}

func InMemoryStorage() StableStorage {
	return StableStorage{
		CurrentTerm:  stable.NewInMemoryIntStore(),
		VotedFor:     stable.NewInMemoryNullableIntStore(),
		Log:          stable.NewInMemoryLogEntriesStore(),
		CommitLength: stable.NewInMemoryIntStore(),
	}
}

type Node struct {
	nodeId   int   // identity of the current node.
	numNodes int   // total number of raft nodes
	nodes    []int // identity of all nodes, including self.

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
	logResponses      chan LogResponse
	logRequests       chan LogRequest
	voteResponses     chan VoteResponse
	voteRequests      chan VoteRequest
	broadcastRequests chan BroadcastRequest

	// On timeout: Candidate runs for election.
	electionTimer   *time.Timer
	electionTimeout time.Duration

	// On timeout: Follower runs for election.
	heartbeatTimer   *time.Timer
	heartbeatTimeout time.Duration

	// On timeout: Leader broadcasts heartbeat.
	replicateLogTicker   *time.Ticker
	replicateLogInterval time.Duration

	// Non-raft stuff
	httpClient *http.Client
	done       chan struct{}
	logger     logrus.FieldLogger
}

func NewNode(nodeId int, storage StableStorage) *Node {
	n := &Node{}
	n.numNodes = 3
	n.nodeId = nodeId
	n.nodes = []int{0, 1, 2}
	n.logger = logrus.New().WithField("nodeId", nodeId)

	n.currentTerm = storage.CurrentTerm
	n.votedFor = storage.VotedFor
	n.log = storage.Log
	n.commitLength = storage.CommitLength

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

	n.logResponses = make(chan LogResponse)
	n.logRequests = make(chan LogRequest)
	n.voteResponses = make(chan VoteResponse)
	n.voteRequests = make(chan VoteRequest)
	n.broadcastRequests = make(chan BroadcastRequest, 10)

}

func (n *Node) initTimers() {
	// Initially disabled
	n.electionTimeout = randomInterval()
	n.electionTimer = &time.Timer{}

	// Initially enabled, and will trigger follower to become candidate.
	n.heartbeatTimeout = randomInterval()
	n.heartbeatTimer = time.NewTimer(n.heartbeatTimeout)

	// Initially disabled
	n.replicateLogInterval = randomInterval()
	n.replicateLogTicker = time.NewTicker(n.replicateLogInterval)
}

func (n *Node) Start() {
	go n.serveRPC()
	go n.mainLoop()
}

func (n *Node) Broadcast(msg []byte) {
	msgReq := BroadcastRequest{
		Msg: msg,
	}
	n.broadcastRequests <- msgReq
}

func (n *Node) Close() {
	close(n.done)
	time.Sleep(50 * time.Millisecond) // todo wait for stuff to close
}
