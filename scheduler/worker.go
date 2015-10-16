/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scheduler

import (
	"time"

	"github.com/pborman/uuid"
)

var workerKillChan = make(chan struct{})

type worker struct {
	id       string
	rcv      <-chan job
	kamikaze chan struct{}
}

func newWorker(rChan <-chan job) *worker {
	return &worker{
		rcv:      rChan,
		id:       uuid.New(),
		kamikaze: make(chan struct{}),
	}
}

// begin a worker
func (w *worker) start() {
	for {
		select {
		case j := <-w.rcv:
			// assert that deadline is not exceeded
			if time.Now().Before(j.Deadline()) {
				j.Run()
				continue
			}
			// reply immediately -- Job not run
			j.ReplChan() <- struct{}{}

		// the single kill-channel -- used when resizing worker pools
		case <-w.kamikaze:
			return

		//the broadcast that kills all workers
		case <-workerKillChan:
			return

		}
	}
}
