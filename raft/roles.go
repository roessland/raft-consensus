package raft

type Role string

const (
	roleFollower  Role = "FOLLOWER"
	roleLeader    Role = "LEADER"
	roleCandidate Role = "CANDIDATE"
)
