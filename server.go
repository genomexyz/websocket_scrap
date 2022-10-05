// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func generate_realtime(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		/*mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)*/
		rand.Seed(time.Now().UnixNano())
		min := 10
		max := 50
		random_num := []byte(strconv.Itoa(rand.Intn(max-min+1) + min))
		//fmt.Println("cek mt", mt)
		//err = c.WriteMessage(mt, message)
		fmt.Println("send data", random_num)
		err = c.WriteMessage(1, random_num)
		if err != nil {
			log.Println("write:", err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("templates", "live_data.html")
	homeTemplate, err := template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	homeTemplate.Execute(w, "ws://"+r.Host+"/real_time")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/real_time", generate_realtime)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
