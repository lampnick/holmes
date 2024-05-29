/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	mlog "mosn.io/pkg/log"
	"net/http"
	"sync"
	"time"

	"github.com/lampnick/holmes"
)

// run `curl http://localhost:10003/lockorder1` after 15s(warn up)
// run `curl http://localhost:10003/lockorder2` after 15s(warn up)
// run `curl http://localhost:10003/req` after 15s(warn up)
func init() {
	http.HandleFunc("/lockorder1", lockorder1)
	http.HandleFunc("/lockorder2", lockorder2)
	http.HandleFunc("/req", req)
	go http.ListenAndServe(":10003", nil) //nolint:errcheck
}

func main() {
	h, _ := holmes.New(
		holmes.WithCollectInterval("5s"),
		holmes.WithDumpPath("./tmp"),
		holmes.WithLogger(holmes.NewFileLog("./tmp/holmes.log", mlog.DEBUG)),
		holmes.WithTextDump(),
		holmes.WithGoroutineDump(10, 25, 2000, 10000, time.Minute),
	)
	h.EnableGoroutineDump().Start()
	time.Sleep(time.Hour)
}

var l1 sync.Mutex
var l2 sync.Mutex

func req(wr http.ResponseWriter, req *http.Request) {
	l1.Lock()
	defer l1.Unlock()
}

func lockorder1(wr http.ResponseWriter, req *http.Request) {
	l1.Lock()
	defer l1.Unlock()

	time.Sleep(time.Minute)

	l2.Lock()
	defer l2.Unlock()
}

func lockorder2(wr http.ResponseWriter, req *http.Request) {
	l2.Lock()
	defer l2.Unlock()

	time.Sleep(time.Minute)

	l1.Lock()
	defer l1.Unlock()
}
