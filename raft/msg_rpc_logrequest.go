package raft

import "github.com/roessland/raft-consensus/raft/raftlog"

// https://youtu.be/uXEYuDwm7e4?t=1255

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

}
