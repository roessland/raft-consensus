package raft

import "fmt"

func (n *Node) mainLoop() {
	for {
		select {

		case <-n.heartbeatTimer.C:
			n.runForElection()

		case <-n.electionTimer.C:
			n.runForElection()

		case voteRequest := <-n.voteRequests:
			n.onReceivingVoteRequest(voteRequest)

		case voteResponse := <-n.voteResponses:
			n.onReceivingVoteResponse(voteResponse)

		// TODO 4: on request to broadcast msg

		case <-n.replicateLogTicker.C:
			// TODO 4: on replicateLogTimeout
			fmt.Println("replicate log ticker triggered")
			if n.currentRole == roleLeader {
				for _, followerID := range n.nodes {
					if followerID != n.nodeId {
						n.replicateLog(n.nodeId, followerID)
					}
				}
			}

		case logRequest := <-n.logRequests:
			n.onReceivingLogRequest(logRequest)

		case logResponse := <-n.logResponses:
			n.onReceivingLogResponse(logResponse)
		}
	}
}
