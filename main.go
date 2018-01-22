package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type realms struct {
	Servers []server `json:"realms"`
}

type server struct {
	Type            string `json:"type"`
	Population      string `json:"population"`
	Queue           bool   `json:"queue"`
	Status          bool   `json:"status"`
	Name            string `json:"name"`
	Slug            string `json:"-"`
	BattleGroup     string `json:"-"`
	Local           string `json:"-"`
	Timezone        string `json:"-"`
	ConnectedRealms string `json:"-"`
}

var realm = flag.String("realm", "medivh", "the individual realm's status to check")
var rURL = ""

func main() {
	cfg := LoadConfig()
	rURL = cfg.APIKey
	flag.Parse()
	r := *realm
	getServerStatus(strings.ToLower(r))
}

func getAllServers() []server {
	// realms := realms{}
	cl := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, rURL, nil)
	if err != nil {
		fmt.Println("Request Error:", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "wow_realm_status_go")

	res, err := cl.Do(req)
	if err != nil {
		fmt.Println("Get Error:", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	var r realms

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Println(err)
	}
	return r.Servers
}

func getServerStatus(name string) {
	sm := serverMap()
	if s, ok := sm[name]; ok {
		fmt.Printf("%s's status is: %s!", strings.Title(name), s)
	} else {
		fmt.Printf("%s not found in list of servers.", strings.Title(name))
	}
}

func serverMap() map[string]string {
	servers := getAllServers()
	sm := make(map[string]string)
	for _, server := range servers {
		sm[strings.ToLower(server.Name)] = boolToStatusString(server.Status)
	}
	return sm
}

func boolToStatusString(status bool) string {
	if status == true {
		return "Up"
	}
	return "Down"
}
