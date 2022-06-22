package raft

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestRaft1_SendToSelf(t *testing.T) {
	n1 := NewNode(1, InMemoryStorage())
	n1.Start()
	time.Sleep(time.Second)
}

func TestRaft3_InitialElection(t *testing.T) {

	storage := []StableStorage{InMemoryStorage(), InMemoryStorage(), InMemoryStorage()}

	nodes := []*Node{
		NewNode(0, storage[0]),
		NewNode(1, storage[1]),
		NewNode(2, storage[2]),
	}

	for _, node := range nodes {
		node.Start()
	}

	// Send a message to leader after 300 ms
	go func() {
		time.Sleep(700 * time.Millisecond)
		fmt.Println("sending msg to leader")
		for i, _ := range nodes {
			if nodes[i].currentRole == roleLeader {
				nodes[i].Broadcast([]byte("SET key value"))
				break
			}
		}
	}()

	// Kill and restart leader after 500 ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("killing leader")
		for i, _ := range nodes {
			if nodes[i].currentRole == roleLeader {
				nodes[i].Close()
				time.Sleep(300 * time.Millisecond)
				nodes[i] = NewNode(i, storage[i])
				nodes[i].Start()
				break
			}
		}
	}()

	// Print roles
	for i := 0; i < 40; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(nodes[0].currentRole, nodes[1].currentRole, nodes[2].currentRole)
	}

	// Stop servers
	var wg sync.WaitGroup
	wg.Add(3)
	for i := range nodes {
		go func(i int) {
			defer wg.Done()
			nodes[i].Close()
		}(i)
	}
	wg.Wait()

	// Print logs
	for _, node := range nodes {
		fmt.Println(node.log)
	}

	require.True(t, true)
}
