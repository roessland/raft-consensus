package raft

import (
	"github.com/roessland/raft-consensus/raft/raftlog"
)

// https://youtu.be/uXEYuDwm7e4?t=1598

// Update follower's logs

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (n *Node) appendEntries(prefixLen int, leaderCommit int, suffix []raftlog.Entry) {
	if len(suffix) > 0 && n.log.Len() > prefixLen {
		index := min(n.log.Len(), prefixLen+len(suffix)) - 1
		if n.log.At(index).Term != suffix[index-prefixLen].Term {
			n.log.Truncate(prefixLen)
		}
	}
	if prefixLen+len(suffix) > n.log.Len() {
		for i := n.log.Len() - prefixLen; i < len(suffix); i++ {
			n.log.Append(suffix[i])
		}
	}
	if leaderCommit > n.commitLength.Get() {
		for i := n.commitLength.Get(); i < leaderCommit; i++ {
			// "deliver log[i] msg to the application"
		}
		n.commitLength.Set(leaderCommit)
	}
}
