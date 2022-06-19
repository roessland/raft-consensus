package raft

type Role int

const (
	roleFollower Role = iota
	roleLeader
	roleCandidate
)
