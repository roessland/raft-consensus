package raft

// https://youtu.be/uXEYuDwm7e4?t=1255

// Leader receiving log acknowledgements

type LogResponseParams struct {
	Follower int
	Term     int
	Ack      int
	Success  bool
}

func (n *Node) onReceivingLogResponse(msg LogResponseParams) {

}
