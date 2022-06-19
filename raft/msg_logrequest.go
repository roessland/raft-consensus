package raft

// https://youtu.be/uXEYuDwm7e4?t=1255

type LogRequestParams struct {
	LeaderID     int
	Term         int
	PrefixLen    int
	PrefixTerm   int
	LeaderCommit int
	Suffix       int
}

func (n *Node) onReceivingLogRequest(msg LogRequestParams) {

}
