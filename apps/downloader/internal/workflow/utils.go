/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package workflow

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
)

// Convert a number to a string and pad with 0s if less than 10
func leftPad(i int) string {
	var x string

	if i < 10 {
		x += "0"
	}

	x += strconv.Itoa(i)

	return x
}

// Strip any non-alphanumeric characters from the string
func removeNonAlnum(s string) string {
	var result strings.Builder
	for i := range len(s) {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func startHeartbeat(ctx context.Context) chan bool {
	logger := activity.GetLogger(ctx)

	logger.Debug("Starting heatbeat")

	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				logger.Debug("Heatbeat stopped")
				return
			default:
				// Heartbeats count towards the Cloud message cost
				logger.Debug("Pausing heartbeat")
				time.Sleep(time.Minute * 2)

				logger.Debug("Sending heatbeat")
				activity.RecordHeartbeat(ctx)
			}
		}
	}()

	return quit
}
