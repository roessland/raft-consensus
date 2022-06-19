package raft

import (
	"math"
	"time"
)

// https://youtu.be/uXEYuDwm7e4?t=735

type VoteResponse struct {
	Type    MsgType
	VoterID int
	Term    int
	Granted bool
}

func (n *Node) onReceivingVoteResponse(msg VoteResponse) {
	n.logger.Printf("%#+v", msg)
	voterId := msg.VoterID
	term := msg.Term
	granted := msg.Granted

	if n.currentRole == roleCandidate && term == n.currentTerm.Get() && granted {
		n.votesReceived.Add(voterId)
	}
	if len(n.votesReceived) >= int(math.Ceil((float64(n.numNodes)+1.0)/2.0)) {
		n.currentRole = roleLeader
		n.currentLeader = &n.nodeId
		n.electionTimer = &time.Timer{}
		for _, follower := range n.nodes {
			if follower == n.nodeId {
				continue
			}
			n.sentLength[follower] = n.log.Len()
			n.ackedLength[follower] = 0
			n.replicateLog(n.nodeId, follower)
		}
	} else if term > n.currentTerm.Get() {
		n.currentTerm.Set(term)
		n.currentRole = roleFollower
		n.votedFor.SetNull()
		n.electionTimer = &time.Timer{}
	}
}
