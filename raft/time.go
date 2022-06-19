package raft

import (
	"math/rand"
	"time"
)

func randomInterval() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}
