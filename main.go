package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	realms := realms{}
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Read Error:", err)
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(body, &realms)
	if jsonErr != nil {
		fmt.Println("JSON Error:", jsonErr)
	}
	s := realms.Servers
	// for _, s := range realms.Servers {
	// 	fmt.Println(s.Name)
	// }
	return s
}

func getServerStatus(name string) {
	servers := getAllServers()
	for _, server := range servers {
		status := boolToStatusString(server.Status)
		// fmt.Printf("Server: %s - Status: %s\n", server.Name, status)
		// fmt.Printf("%t\n", name == strings.ToLower(server.Name))
		if server.Name == strings.ToLower(name) {
			fmt.Printf("%s is %s", name, status)
			break
		}
	}
}

func boolToStatusString(status bool) string {
	if status == true {
		return "up"
	}
	return "down"
}
