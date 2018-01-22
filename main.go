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

const rURL = "https://us.api.battle.net/wow/realm/status?locale=en_US&apikey=m8m9as776592afjxmwx45vy8yabgpngb"

var realm = flag.String("realm", "medivh", "the individual realm's status to check")

func main() {
	flag.Parse()
	r := *realm
	getServerStatus(r)
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
	servers := getAllServers()
	for _, server := range servers {
		status := boolToStatusString(server.Status)
		if strings.ToLower(server.Name) == strings.ToLower(name) {
			fmt.Printf("%s's status is: %s!", strings.Title(name), status)
		}
	}
}

func boolToStatusString(status bool) string {
	if status == true {
		return "Up"
	}
	return "Down"
}
