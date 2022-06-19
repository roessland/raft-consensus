package raft

// https://youtu.be/uXEYuDwm7e4?t=735

type VoteResponseParams struct {
	VotedID int
	Term    int
	Granted bool
}

func (n *Node) onReceivingVoteResponse(msg VoteResponseParams) {

}
