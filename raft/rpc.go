package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type MsgType string

const (
	MsgTypeLogRequest   = "LogRequest"
	MsgTypeLogResponse  = "LogResponse"
	MsgTypeVoteRequest  = "VoteRequest"
	MsgTypeVoteResponse = "VoteResponse"
)

func (n *Node) httpMsgHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	checkErr(err)

	var someMsg genericMsg
	err = json.Unmarshal(body, &someMsg)
	checkErr(err)

	switch someMsg.Type {
	case MsgTypeLogResponse:
		msg := LogResponse{}
		err = json.Unmarshal(body, &msg)
		checkErr(err)
		n.logResponses <- msg
	case MsgTypeLogRequest:
		msg := LogRequest{}
		err = json.Unmarshal(body, &msg)
		checkErr(err)
		n.logRequests <- msg
	case MsgTypeVoteResponse:
		msg := VoteResponse{}
		err = json.Unmarshal(body, &msg)
		checkErr(err)
		n.voteResponses <- msg
	case MsgTypeVoteRequest:
		msg := VoteRequest{}
		err = json.Unmarshal(body, &msg)
		checkErr(err)
		n.voteRequests <- msg
	default:
		log.Printf("got RPC message of unknown type: %s", someMsg.Type)
	}

	n.logger.Infof("%d got RPC: %s", n.nodeId, string(body))
}

func (n *Node) serveRPC() {
	r := mux.NewRouter()
	r.HandleFunc("/", n.httpMsgHandler)
	addr := fmt.Sprintf("127.0.0.1:%d", 50000+n.nodeId)
	n.logger.Infof("Raft: listening on HTTP at %s", addr)
	server := &http.Server{Addr: addr, Handler: r}
	go func() {
		<-n.done
		n.logger.Infof("closing rpc server")
		checkErr(server.Close())
	}()
	n.logger.Infof("rpc server: %s", server.ListenAndServe())
}

// sendRPC sends a message.
func (n *Node) sendRPC(dstNodeId int, msg any) {
	url := fmt.Sprintf("http://127.0.0.1:%d/", 50000+dstNodeId)
	body, err := json.Marshal(&msg)
	checkErr(err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	checkErr(err)

	go func() {
		resp, err := n.httpClient.Do(req)
		if err != nil {
			printErr(err)
			return
		}

		if resp.StatusCode != 200 {
			n.logger.Infof("non-200-status")
		}
	}()
}
