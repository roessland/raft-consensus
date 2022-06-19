package raft

// https://youtu.be/uXEYuDwm7e4?t=521

type VoteRequestParams struct {
	CandidateID        int
	CandidateTerm      int
	CandidateLogLength int
	CandidateLogTerm   int
}

func (n *Node) onReceivingVoteRequest(msg VoteRequestParams) {

}
