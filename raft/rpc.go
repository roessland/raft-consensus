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

func (n *Node) httpMsgHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	checkErr(err)
	log.Print("got RPC:", string(body))
}

func (n *Node) serveRPC() {
	r := mux.NewRouter()
	r.HandleFunc("/", n.httpMsgHandler)
	addr := fmt.Sprintf("127.0.0.1:%d", 50000+n.nodeId)
	log.Printf("Raft: listening on HTTP at %s", addr)
	server := &http.Server{Addr: addr, Handler: r}
	log.Fatal(server.ListenAndServe())
}

func (n *Node) sendRPC(dstNodeId int, msg any) {
	url := fmt.Sprintf("http://127.0.0.1:%d/", 50000+dstNodeId)
	body, err := json.Marshal(&msg)
	checkErr(err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	checkErr(err)

	resp, err := n.httpClient.Do(req)
	checkErr(err)

	if resp.StatusCode != 200 {
		log.Println("non-200-status")
	}
}
