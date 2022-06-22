package raft

import (
	"github.com/roessland/raft-consensus/raft/sets"
	"time"
)

func (n *Node) runForElection() {
	n.logger.Info("running for election")
	n.currentTerm.Set(n.currentTerm.Get() + 1)
	n.currentRole = roleCandidate
	n.votedFor.Set(n.nodeId)
	n.votesReceived = sets.NewIntSet(n.nodeId)
	lastTerm := 0

	if n.log.Len() > 0 {
		lastTerm = n.log.At(n.log.Len() - 1).Term
	}

	msg := VoteRequest{
		Type:               MsgTypeVoteRequest,
		CandidateID:        n.nodeId,
		CandidateTerm:      n.currentTerm.Get(),
		CandidateLogLength: n.log.Len(),
		CandidateLogTerm:   lastTerm,
	}

	for _, node := range n.nodes {
		n.sendRPC(node, msg)
	}

	n.electionTimer = time.NewTimer(n.electionTimeout)
}
