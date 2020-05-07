package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/laser/arctic-logs-webservice/types"
)

var alphaSortedClanNames = make([]string, 0)
var alphaSortedLogNames = make([]string, 0)
var alphaSortedPlayerNames = make([]string, 0)
var logs = make(map[string]string)
var logsByClanName = make(map[string][]string)
var logsByPlayerName = make(map[string][]string)

func init() {
	var metadataFiles []string

	err := filepath.Walk("./logs", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".meta" {
			return nil
		}

		metadataFiles = append(metadataFiles, path)

		return nil
	})
	check(err)

	for _, file := range metadataFiles {
		buf, err := ioutil.ReadFile(file)
		check(err)

		m := types.Meta{}
		err = json.Unmarshal(buf, &m)
		check(err)

		orig := strings.Split(file, ".meta")[0]
		orig = filepath.Base(orig)

		for idx := range m.PlayerNames {
			if _, ok := logsByPlayerName[m.PlayerNames[idx]]; !ok {
				logsByPlayerName[m.PlayerNames[idx]] = []string{orig}
			} else {
				logsByPlayerName[m.PlayerNames[idx]] = append(logsByPlayerName[m.PlayerNames[idx]], orig)
			}
		}

		for idx := range m.ClanNames {
			if _, ok := logsByClanName[m.ClanNames[idx]]; !ok {
				logsByClanName[m.ClanNames[idx]] = []string{orig}
			} else {
				logsByClanName[m.ClanNames[idx]] = append(logsByClanName[m.ClanNames[idx]], orig)
			}
		}

		logs[orig] = orig
	}

	for idx := range logsByClanName {
		alphaSortedClanNames = append(alphaSortedClanNames, idx)
	}

	for idx := range logsByPlayerName {
		alphaSortedPlayerNames = append(alphaSortedPlayerNames, idx)
	}

	for idx := range logs {
		alphaSortedLogNames = append(alphaSortedLogNames, idx)
	}

	sort.Strings(alphaSortedClanNames)
	sort.Strings(alphaSortedPlayerNames)
	sort.Strings(alphaSortedLogNames)
}

var ClanDetailPageHandler = createDetailPageHandler(logsByClanName)

var PlayerDetailPageHandler = createDetailPageHandler(logsByPlayerName)

var LogPageHandler = func(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	f, err := os.Open(fmt.Sprintf("./logs/%s", name))
	defer f.Close()
	check(err)

	_, err = io.Copy(w, f)
	check(err)
}

var HomePageHandler = func(w http.ResponseWriter, r *http.Request) {
	tpl := "home.tpl"

	data := types.HomePage{
		Players: []types.Link{},
		Clans:   []types.Link{},
		Logs:    []types.Link{},
	}

	for _, v := range alphaSortedClanNames {
		data.Clans = append(data.Clans, types.Link{
			Label: v,
			Url:   fmt.Sprintf("/clans/%s", v),
		})
	}

	for _, v := range alphaSortedPlayerNames {
		data.Players = append(data.Players, types.Link{
			Label: v,
			Url:   fmt.Sprintf("/players/%s", v),
		})
	}

	for _, v := range alphaSortedLogNames {
		data.Logs = append(data.Logs, types.Link{
			Label: v,
			Url:   fmt.Sprintf("/logs/%s", logs[v]),
		})
	}

	t, err := template.ParseFiles(tpl)
	check(err)

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	check(err)

	_, err = io.Copy(w, &buf)
	check(err)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}

	rtr := mux.NewRouter()

	rtr.HandleFunc("/clans/{name}", ClanDetailPageHandler).Methods("GET")
	rtr.HandleFunc("/players/{name}", PlayerDetailPageHandler).Methods("GET")
	rtr.HandleFunc("/logs/{name}", LogPageHandler).Methods("GET")
	rtr.HandleFunc("/", HomePageHandler).Methods("GET")

	http.Handle("/", rtr)

	fmt.Printf("server listening on port %s...\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func createDetailPageHandler(idx map[string][]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl := "detail.tpl"

		name := mux.Vars(r)["name"]

		logs := idx[name]

		data := types.DetailPage{
			Logs: []types.Link{},
		}

		for idx := range logs {
			data.Logs = append(data.Logs, types.Link{
				Label: logs[idx],
				Url:   fmt.Sprintf("/logs/%s", logs[idx]),
			})
		}

		t, err := template.ParseFiles(tpl)
		check(err)

		var buf bytes.Buffer
		err = t.Execute(&buf, data)
		check(err)

		_, err = io.Copy(w, &buf)
		check(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
