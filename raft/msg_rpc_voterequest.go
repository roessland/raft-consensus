package raft

// https://youtu.be/uXEYuDwm7e4?t=521

type VoteRequest struct {
	Type               MsgType
	CandidateID        int
	CandidateTerm      int
	CandidateLogLength int
	CandidateLogTerm   int
}

func (n *Node) onReceivingVoteRequest(msg VoteRequest) {
	if msg.CandidateTerm > n.currentTerm.Get() {
		n.currentTerm.Set(msg.CandidateTerm)
		n.currentRole = roleFollower
		n.votedFor.SetNull()
	}
	lastTerm := 0
	if n.log.Len() > 0 {
		lastTerm = n.log.At(n.log.Len() - 1).Term
	}
	logOk := (msg.CandidateLogTerm > lastTerm) ||
		(msg.CandidateLogTerm == lastTerm && msg.CandidateLogLength >= n.log.Len())

	voteResponseMsg := VoteResponse{
		Type:    MsgTypeVoteResponse,
		VoterID: n.nodeId,
		Term:    n.currentTerm.Get(),
	}
	votedForCandidateOrNobody := n.votedFor.IsNull() || n.votedFor.Get() == msg.CandidateID
	if msg.CandidateTerm == n.currentTerm.Get() && logOk && votedForCandidateOrNobody {
		n.votedFor.Set(msg.CandidateID)
		voteResponseMsg.Granted = true
		n.sendRPC(msg.CandidateID, voteResponseMsg)
	} else {
		voteResponseMsg.Granted = false
		n.sendRPC(msg.CandidateID, voteResponseMsg)
	}
}
