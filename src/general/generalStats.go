/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

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
package general

import (
	"context"
	"sync"
	"time"

	"github.com/pesos/grofer/src/utils"
)

// GlobalStats gets stats about the mem and the CPUs and prints it.
func GlobalStats(ctx context.Context,
	dataChannel chan utils.DataStats,
	refreshRate int32) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default: // Get Memory and CPU rates per core periodically
			wg := sync.WaitGroup{}
			errCh := make(chan error, 4)

			wg.Add(4)
			go func(dc chan utils.DataStats) {
				defer wg.Done()
				errCh <- ServeCPURates(dc)
			}(dataChannel)

			go func(dc chan utils.DataStats) {
				defer wg.Done()
				errCh <- ServeMemRates(dc)
			}(dataChannel)

			go func(dc chan utils.DataStats) {
				defer wg.Done()
				errCh <- ServeDiskRates(dc)
			}(dataChannel)

			go func(dc chan utils.DataStats) {
				defer wg.Done()
				errCh <- ServeNetRates(dc)
			}(dataChannel)

			wg.Wait()
			close(errCh)

			for err := range errCh {
				if err != nil {
					return err
				}
			}

			time.Sleep(time.Duration(refreshRate) * time.Millisecond)
		}
	}
}
