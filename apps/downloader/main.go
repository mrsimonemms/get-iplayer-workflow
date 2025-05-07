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

package main

import (
	"log/slog"

	"github.com/mrsimonemms/get-iplayer-workflow/apps/downloader/internal/config"
	"github.com/mrsimonemms/get-iplayer-workflow/apps/downloader/internal/workflow"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"go.temporal.io/sdk/client"
	tLog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to load config")
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Str("logLevel", cfg.LogLevel).Msg("Unknown log level")
	}
	zerolog.SetGlobalLevel(logLevel)

	host := cfg.Temporal.Address
	namespace := cfg.Temporal.Namespace

	log.Debug().Msg("Connecting to NATS server")
	natsClient, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to NATS server")
	}
	defer natsClient.Close()

	log.Debug().Msg("Connecting to Temporal server")
	temporalClient, err := client.Dial(client.Options{
		HostPort:  host,
		Namespace: namespace,
		Logger: tLog.NewStructuredLogger(slog.New(slogzerolog.Option{
			Logger: &log.Logger,
		}.NewZerologHandler())),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create Temporal client")
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, "downloadByPID", worker.Options{})

	w.RegisterWorkflow(workflow.DownloadBBCProgramme)
	w.RegisterWorkflow(workflow.ParseDownloadedProgramme)

	activities := workflow.NewActivities(natsClient, cfg)
	w.RegisterActivity(activities)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatal().Err(err).Msg("Unable to start worker")
	}
}
