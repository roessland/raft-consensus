package raft

import (
	"math"
)

// https://youtu.be/uXEYuDwm7e4?t=2136

// Leader committing log entries

func (n *Node) commitLogEntries() {
	for n.commitLength.Get() < n.log.Len() {
		acks := 0
		for node := range n.nodes {
			if n.ackedLength[node] > n.commitLength.Get() {
				acks = acks + 1
			}
		}
		if acks >= int(math.Ceil(float64(n.numNodes+1)/2)) {
			// TODO "deliver log[commitLength].msg to the application"
			n.logger.Info("deliver to application", n.log.At(n.commitLength.Get()))
			n.commitLength.Set(n.commitLength.Get() + 1)
		} else {
			break
		}
	}
}
