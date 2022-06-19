package raft

import "time"

// https://youtu.be/uXEYuDwm7e4?t=1255

// Leader receiving log acknowledgements

type LogResponse struct {
	Type     MsgType
	Follower int
	Term     int
	Ack      int
	Success  bool
}

func (n *Node) onReceivingLogResponse(msg LogResponse) {
	follower := msg.Follower
	term := msg.Term
	ack := msg.Ack
	success := msg.Success

	if term == n.currentTerm.Get() && n.currentRole == roleLeader {
		if success && ack >= n.ackedLength[follower] {
			n.sentLength[follower] = ack
			n.ackedLength[follower] = ack
			n.commitLogEntries()
		} else if n.sentLength[follower] > 0 {
			n.sentLength[follower] = n.sentLength[follower] - 1
			n.replicateLog(n.nodeId, follower)
		}
	} else if term > n.currentTerm.Get() {
		n.currentTerm.Set(term)
		n.currentRole = roleFollower
		n.votedFor.SetNull()
		n.electionTimer = &time.Timer{}
	}
}
