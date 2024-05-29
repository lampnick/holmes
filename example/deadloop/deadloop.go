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
	"time"

	"github.com/lampnick/holmes"
)

// run `curl http://localhost:10003/alldeadloopoc` after 15s(warn up)
func init() {
	http.HandleFunc("/deadloop", deadloop)
	go http.ListenAndServe(":10003", nil) //nolint:errcheck
}

func main() {
	h, _ := holmes.New(
		holmes.WithCollectInterval("2s"),
		holmes.WithDumpPath("./tmp"),
		holmes.WithLogger(holmes.NewFileLog("./tmp/holmes.log", mlog.DEBUG)),
		holmes.WithCPUDump(10, 25, 80, time.Minute),
	)
	h.EnableCPUDump().Start()
	time.Sleep(time.Hour)
}

func deadloop(wr http.ResponseWriter, req *http.Request) {
	for {
		select {
		case <-req.Context().Done():
			break
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
