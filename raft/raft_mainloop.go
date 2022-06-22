package raft

func (n *Node) mainLoop() {
	for {
		select {

		case <-n.heartbeatTimer.C:
			n.logger.Info("leader timeout triggered")
			n.runForElection()

		case <-n.electionTimer.C:
			n.logger.Info("election timer triggered")
			n.runForElection()

		case voteRequest := <-n.voteRequests:
			n.logger.Info("got vote request")
			n.onReceivingVoteRequest(voteRequest)

		case voteResponse := <-n.voteResponses:
			n.logger.Info("got vote response")
			n.onReceivingVoteResponse(voteResponse)

		// TODO 4: on request to broadcast msg

		case <-n.replicateLogTicker.C:
			// TODO 4: on replicateLogTimeout
			n.logger.Info("replicate log ticker triggered")
			if n.currentRole == roleLeader {
				for _, followerID := range n.nodes {
					if followerID != n.nodeId {
						n.replicateLog(n.nodeId, followerID)
					}
				}
			}

		case logRequest := <-n.logRequests:
			n.logger.Info("got log request")
			n.onReceivingLogRequest(logRequest)

		case logResponse := <-n.logResponses:
			n.logger.Info("got log response")
			n.onReceivingLogResponse(logResponse)
		}
	}
}
