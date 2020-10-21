package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/vitecoin/zi/api"
)

type request struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Line  string `json:"line"`
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	K, ok := r.URL.Query()["key"]
	if ok != true {
		w.Write([]byte("key not found"))
	} else {
		parsed := api.Init()
		data := api.Get(parsed, K[0], true)
		Json, err := json.Marshal(request{Key: data.Key, Value: data.Value, Line: strconv.Itoa(data.Line)})
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(string(Json)))
	}

}

type SetPair struct {
	Value string `json:"Value"`
	Key   string `json:"Key"`
}

func set(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	K, ok := r.URL.Query()["data"]
	if ok != true {
		w.Write([]byte("Data not found"))
	} else {
		s := string(K[0])
		data := SetPair{}
		json.Unmarshal([]byte(s), &data)
		api.Set(api.Pair{Key: data.Key, Value: data.Value}, true)
		parsed := api.Init()
		get := api.Get(parsed, data.Key, false)
		json, err := json.Marshal(api.Pair{Line: get.Line, Value: get.Value, Key: get.Key})
		if err != nil {
			panic(err)
		}
		w.Write([]byte(string(json)))
	}
}
func dump(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	K, ok := r.URL.Query()["data"]
	P, okP := r.URL.Query()["path"]
	if ok != true || okP != true {
		w.Write([]byte("Data or path not found"))
	} else {
		s := "+" + string(K[0])
		data := SetPair{}
		json.Unmarshal([]byte(s), &data)
		api.Dump(data.Key, data.Value, P[0], true)
		w.Write([]byte("OK"))
	}
}
func del(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	K, ok := r.URL.Query()["key"]
	if ok != true {
		w.Write([]byte("Key not found"))
	} else {
		api.Del(K[0], true)
		w.Write([]byte("OK"))
	}
}
func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write([]byte(api.GetAll()))
}
func bind(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	url, okUrl := r.URL.Query()["url"]
	key, okKey := r.URL.Query()["key"]
	if okUrl == true && okKey == true {
		api.Bind(key[0], url[0], true)
		w.Write([]byte("OK"))
	} else {
		w.Write([]byte("Key or url not found"))
	}
}
func rename(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	new, okUrl := r.URL.Query()["new"]
	origin, okKey := r.URL.Query()["origin"]
	if okUrl == true && okKey == true {
		api.Rename(origin[0], new[0], true)
		w.Write([]byte("OK"))

	} else {
		w.Write([]byte("New or origin not found"))
	}
}
func getrow(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	key, ok := r.URL.Query()["key"]
	if ok == false {
		w.Write([]byte("Key not found"))
	} else {
		r := api.GetRow(api.Init(), key[0])
		data, _ := json.Marshal(r)
		w.Write([]byte(data))
	}
}
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// Serve starts database server.
func Serve(port string) {
	if _, err := os.Stat("dump.zi"); err != nil {
		ioutil.WriteFile("dump.zi", []byte(""), 0644)
	}
	fmt.Println("Server running on port " + port)
	fmt.Printf("Public address: http://%s:%s\n", getOutboundIP(), port)
	fmt.Printf("Localhost address: http://127.0.0.1:%s\n", port)
	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
	http.HandleFunc("/del", del)
	http.HandleFunc("/getrow", getrow)
	http.HandleFunc("/getall", getAll)
	http.HandleFunc("/bind", bind)
	http.HandleFunc("/rename", rename)
	http.HandleFunc("/dump", dump)
	http.HandleFunc("/", root)
	http.ListenAndServe(":"+port, nil)
	// if err != nil {
	// log.Fatal(err)
	// }
}
