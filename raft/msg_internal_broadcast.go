package raft

import (
	"github.com/roessland/raft-consensus/raft/raftlog"
	"log"
)

// https://youtu.be/uXEYuDwm7e4?t=1059

type BroadcastRequest struct {
	Msg []byte
}

func (n *Node) onReceivingBroadcastRequest(msgReq BroadcastRequest) {
	msg := msgReq.Msg

	if n.currentRole == roleLeader {
		entry := raftlog.Entry{
			Term: n.currentTerm.Get(),
			Msg:  msg,
		}
		n.log.Append(entry)
		n.ackedLength[n.nodeId] = n.log.Len()
		for follower := range n.nodes {
			if follower == n.nodeId {
				continue
			}
			n.replicateLog(n.nodeId, follower)
		}
	} else {
		// forward the request to currentLeader via a FIFO link.
		log.Printf("i'm not the leader. Ask %d instead", n.currentLeader)
	}
}
