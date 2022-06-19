package stable

import "github.com/roessland/raft-consensus/raft"

// In-memory implementations of stable storage,
// used for testing.

// InMemoryIntStore is a fake IntStore used for testing.
type InMemoryIntStore struct {
	val    int
	exists bool
}

func NewInMemoryIntStore() IntStore {
	return &InMemoryIntStore{}
}

func (s *InMemoryIntStore) Get() int {
	return s.val
}

func (s *InMemoryIntStore) Set(val int) {
	s.exists = true
	s.val = val
}

func (s *InMemoryIntStore) AlreadyExists() bool {
	return s.exists
}

// InMemoryNullableIntStore is a fake NullableIntStore used for testing.
type InMemoryNullableIntStore struct {
	InMemoryIntStore
}

func NewInMemoryNullableIntStore() NullableIntStore {
	return &InMemoryNullableIntStore{}
}

func (s *InMemoryNullableIntStore) Get() int {
	if s.val == -1 {
		panic("got null value")
	}
	return s.val
}

func (s *InMemoryNullableIntStore) Set(val int) {
	if val == -1 {
		panic("use SetNull")
	}
	s.val = val
}

func (s *InMemoryNullableIntStore) AlreadyExists() bool {
	return s.InMemoryIntStore.AlreadyExists()
}

func (s *InMemoryNullableIntStore) IsNull() bool {
	return s.val == -1
}

func (s *InMemoryNullableIntStore) SetNull() {
	s.val = -1
}

// InMemoryLogEntriesStore is a fake LogEntriesStore used for testing.
type InMemoryLogEntriesStore struct {
	Entries []raft.LogEntry
	Exists  bool
}

func NewInMemoryLogEntriesStore() LogEntriesStore {
	return &InMemoryLogEntriesStore{}
}

func (s *InMemoryLogEntriesStore) SetEmpty() {
	s.Exists = true
}

func (s *InMemoryLogEntriesStore) Len() int {
	return len(s.Entries)
}

func (s *InMemoryLogEntriesStore) At(idx int) raft.LogEntry {
	return s.Entries[idx]
}

func (s *InMemoryLogEntriesStore) Append(entry raft.LogEntry) {
	s.Exists = true
	s.Entries = append(s.Entries, entry)
}

func (s *InMemoryLogEntriesStore) AlreadyExists() bool {
	return s.Exists
}
