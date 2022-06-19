package raft

import "github.com/roessland/raft-consensus/raft/raftlog"

// https://youtu.be/uXEYuDwm7e4?t=1165

// replicateLog is called on the leader whenever there is a new message in the log,
// and also periodically. If there are no new messages, suffix
// is the empty list. LogRequest messages with suffix=[] serve as
// heartbeats, letting followers know that the leader is still alive.
func (n *Node) replicateLog(leaderID int, followerID int) {
	// TODO https://youtu.be/uXEYuDwm7e4?t=1165
	prefixLen := n.sentLength[followerID]
	var suffix []raftlog.Entry
	for i := prefixLen; i < n.log.Len(); i++ {
		suffix = append(suffix, n.log.At(i))
	}
	prefixTerm := 0
	if prefixLen > 0 {
		prefixTerm = n.log.At(prefixLen - 1).Term
	}

	msg := LogRequest{
		Type:         MsgTypeLogRequest,
		LeaderID:     leaderID,
		Term:         n.currentTerm.Get(),
		PrefixLen:    prefixLen,
		PrefixTerm:   prefixTerm,
		LeaderCommit: n.commitLength.Get(),
		Suffix:       suffix,
	}
	n.sendRPC(followerID, msg)
}
