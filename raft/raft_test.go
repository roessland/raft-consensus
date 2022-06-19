package raft

import "testing"

func TestRaft1_SendToSelf(t *testing.T) {
	n1 := NewNode(1)
	defer n1.Close()
	n1.sendRPC(1, "hello")
}
