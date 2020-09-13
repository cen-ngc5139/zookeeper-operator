package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/ghostbaby/zk-agent/api"

	"github.com/ghostbaby/zk-agent/g"

	"github.com/samuel/go-zookeeper/zk"
)

type AddMember struct {
	Record string `json:"record"`
}

func getZkStatus(w http.ResponseWriter, r *http.Request) {

	if !g.Config().Http.Backdoor {
		w.Write([]byte("/run disabled"))
		return
	}

	c, _ := zk.FLWSrvr([]string{g.Config().ZkHost}, time.Second*10) //*10)

	var out []byte
	var err error
	for _, v := range c {
		out, err = json.Marshal(v)
		if err != nil {
			w.Write([]byte("exec fail: " + err.Error()))
			return
		}
	}

	w.Write(out)

}

func Health(w http.ResponseWriter, r *http.Request) {

	if err := api.WriteJSON(w, "healthy"); err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}

}

func getZkClient(w http.ResponseWriter, r *http.Request) {
	if !g.Config().Http.Backdoor {
		w.Write([]byte("/run disabled"))
		return
	}

	c, _ := zk.FLWCons([]string{g.Config().ZkHost}, time.Second*10) //*10)

	var out []byte
	var err error
	for _, v := range c {
		out, err = json.Marshal(v)
		if err != nil {
			w.Write([]byte("exec fail: " + err.Error()))
			return
		}
	}

	w.Write(out)

}

func getZkRunok(w http.ResponseWriter, r *http.Request) {

	if !g.Config().Http.Backdoor {
		w.Write([]byte("/run disabled"))
		return
	}

	c := zk.FLWRuok([]string{g.Config().ZkHost}, time.Second*10) //*10)

	var out []byte
	var err error
	for _, v := range c {
		out, err = json.Marshal(v)
		if err != nil {
			w.Write([]byte("exec fail: " + err.Error()))
			return
		}
	}

	w.Write(out)

}

func addMember(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 {
		http.Error(w, "body is blank", http.StatusBadRequest)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	record := &AddMember{}
	if err := json.Unmarshal(bs, record); err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}

	body := record.Record

	zkConnect, _, err := zk.Connect([]string{g.Config().ZkHost}, time.Second*10)
	if err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}
	defer zkConnect.Close()
	var reconfigData []string
	reconfigData = append(reconfigData, body)

	out, err := zkConnect.IncrementalReconfig(reconfigData, nil, -1)
	if err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}
	request, _ := json.Marshal(out)

	w.Write(request)
}

func getMember(w http.ResponseWriter, r *http.Request) {

	var out []byte
	var err error

	if !g.Config().Http.Backdoor {
		w.Write([]byte("/run disabled"))
		return
	}

	zkConnect, _, err := zk.Connect([]string{g.Config().ZkHost}, time.Second*10)
	if err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}
	defer zkConnect.Close()

	out, _, err = zkConnect.Get("/zookeeper/config")
	if err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}

	record := &AddMember{
		Record: string(out),
	}

	recordByte, err := json.Marshal(record)
	if err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}

	w.Write(recordByte)

}

func delMember(w http.ResponseWriter, r *http.Request) {

	if r.ContentLength == 0 {
		http.Error(w, "body is blank", http.StatusBadRequest)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	record := &AddMember{}
	if err := json.Unmarshal(bs, record); err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}

	body := record.Record

	//获取需要删除的myid
	zkConfigRegexp := regexp.MustCompile(`^server.(\d+)=.*2181$`)
	params := zkConfigRegexp.FindStringSubmatch(body)
	for _, param := range params {
		fmt.Println(param)
	}

	if len(params) <= 0 {
		err := errors.New("unable to get zk node id")
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}

	zkNodeID := params[1]

	zkConnect, _, err := zk.Connect([]string{g.Config().ZkHost}, time.Second*10)
	if err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}
	defer zkConnect.Close()

	var reconfigData []string
	reconfigData = append(reconfigData, zkNodeID)

	out, err := zkConnect.IncrementalReconfig(nil, reconfigData, -1)
	if err != nil {
		w.Write([]byte("exec fail: " + err.Error()))
		return
	}
	request, _ := json.Marshal(out)

	w.Write(request)
}
