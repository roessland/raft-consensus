package stable

import "github.com/roessland/raft-consensus/raft"

// Stable variables stored on disk are not lost if the process crashes.
// https://youtu.be/uXEYuDwm7e4?t=213

/*
Some node variables cannot be lost in case of a crash.
They must be written and flushed to disk every time they are modified.
They must be reloaded from disk on crash recovery.
*/

// IntStore is an int that is stored to disk on every write.
// Only the initial Get() reads from disk. Subsequent reads are from memory.
type IntStore interface {
	Get() int
	Set(int)
	AlreadyExists() bool
}

// NullableIntStore is an int that is stored to disk on every write.
// Only the initial Get() reads from disk. Subsequent reads are from memory.
type NullableIntStore interface {
	Get() int
	Set(int)
	AlreadyExists() bool
	IsNull() bool
	SetNull()
}

type LogEntriesStore interface {
	SetEmpty() // SetEmpty stores an empty list.
	Len() int
	At(int) raft.LogEntry
	Append(entry raft.LogEntry)
	AlreadyExists() bool
}
