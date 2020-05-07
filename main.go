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
	"github.com/laser/davidwees-arctic-logs/types"
)

var playerLogs = make(map[string][]string)
var playerLogsSortedKeys = make([]string, 0)

var clanLogs = make(map[string][]string)
var clanLogsSortedKeys = make([]string, 0)

var logs = make(map[string]string)
var logsSortedKeys = make([]string, 0)

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

		for idx := range m.CharNames {
			if _, ok := playerLogs[m.CharNames[idx]]; !ok {
				playerLogs[m.CharNames[idx]] = []string{orig}
			} else {
				playerLogs[m.CharNames[idx]] = append(playerLogs[m.CharNames[idx]], orig)
			}
		}

		for idx := range m.ClanNames {
			if _, ok := clanLogs[m.ClanNames[idx]]; !ok {
				clanLogs[m.ClanNames[idx]] = []string{orig}
			} else {
				clanLogs[m.ClanNames[idx]] = append(clanLogs[m.ClanNames[idx]], orig)
			}
		}

		logs[orig] = orig
	}

	for idx := range clanLogs {
		clanLogsSortedKeys = append(clanLogsSortedKeys, idx)
	}

	for idx := range playerLogs {
		playerLogsSortedKeys = append(playerLogsSortedKeys, idx)
	}

	for idx := range logs {
		logsSortedKeys = append(logsSortedKeys, idx)
	}

	sort.Strings(clanLogsSortedKeys)
	sort.Strings(playerLogsSortedKeys)
	sort.Strings(logsSortedKeys)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}

	rtr := mux.NewRouter()

	rtr.HandleFunc("/clans/{name}", func(w http.ResponseWriter, r *http.Request) {
		tpl := "detail.tpl"

		name := mux.Vars(r)["name"]

		logs := clanLogs[name]

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
	}).Methods("GET")

	rtr.HandleFunc("/players/{name}", func(w http.ResponseWriter, r *http.Request) {
		tpl := "detail.tpl"

		name := mux.Vars(r)["name"]

		logs := playerLogs[name]

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
	}).Methods("GET")

	rtr.HandleFunc("/logs/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		f, err := os.Open(fmt.Sprintf("./logs/%s", name))
		defer f.Close()
		check(err)

		_, err = io.Copy(w, f)
		check(err)
	}).Methods("GET")

	rtr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl := "home.tpl"

		data := types.HomePage{
			Players: []types.Link{},
			Clans:   []types.Link{},
			Logs:    []types.Link{},
		}

		for _, v := range clanLogsSortedKeys {
			data.Clans = append(data.Clans, types.Link{
				Label: v,
				Url:   fmt.Sprintf("/clans/%s", v),
			})
		}

		for _, v := range playerLogsSortedKeys {
			data.Players = append(data.Players, types.Link{
				Label: v,
				Url:   fmt.Sprintf("/players/%s", v),
			})
		}

		for _, v := range logsSortedKeys {
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
	}).Methods("GET")

	http.Handle("/", rtr)

	fmt.Printf("server listening on port %s...\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
