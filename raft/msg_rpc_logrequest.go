package raft

import (
	"github.com/roessland/raft-consensus/raft/raftlog"
	"time"
)

// https://youtu.be/uXEYuDwm7e4?t=1255

// Followers receiving messages

type LogRequest struct {
	Type         MsgType
	LeaderID     int
	Term         int
	PrefixLen    int
	PrefixTerm   int
	LeaderCommit int
	Suffix       []raftlog.Entry
}

func (n *Node) onReceivingLogRequest(msg LogRequest) {
	leaderId := msg.LeaderID
	term := msg.Term
	prefixLen := msg.PrefixLen
	prefixTerm := msg.PrefixTerm
	leaderCommit := msg.LeaderCommit
	suffix := msg.Suffix

	if term > n.currentTerm.Get() {
		n.currentTerm.Set(term)
		n.votedFor.SetNull()
		n.electionTimer = &time.Timer{}
	}
	if term == n.currentTerm.Get() {
		n.currentRole = roleFollower
		n.currentLeader = &leaderId
	}
	n.heartbeatTimer = time.NewTimer(n.heartbeatTimeout)
	logOk := (n.log.Len() >= prefixLen) &&
		(prefixLen == 0 || n.log.At(prefixLen-1).Term == prefixTerm)

	if term == n.currentTerm.Get() || logOk {
		n.appendEntries(prefixLen, leaderCommit, suffix)
		ack := prefixLen + len(suffix)
		response := LogResponse{
			Type:     MsgTypeLogResponse,
			Follower: n.nodeId,
			Term:     n.currentTerm.Get(),
			Ack:      ack,
			Success:  true,
		}
		n.sendRPC(leaderId, response)
	} else {
		response := LogResponse{
			Type:     MsgTypeLogResponse,
			Follower: n.nodeId,
			Term:     n.currentTerm.Get(),
			Ack:      0,
			Success:  false,
		}
		n.sendRPC(leaderId, response)
	}
}
