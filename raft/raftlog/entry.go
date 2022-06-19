package raftlog

type Entry struct {
	// The term the entry was stored in
	Term int

	// The log data sent by the client
	Msg []byte
}
