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

package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/client"
)

type NATS struct {
	URL string `env:"URL"`
}

type Temporal struct {
	Address   string `env:"ADDRESS"`
	Namespace string `env:"NAMESPACE"`
}

type Config struct {
	NATS     `envPrefix:"NATS_"`
	Temporal `envPrefix:"TEMPORAL_"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	OutputDir string `env:"OUTPUT_DIR,required"`
}

func Load() (*Config, error) {
	cfg := Config{
		NATS: NATS{
			URL: nats.DefaultURL,
		},
		Temporal: Temporal{
			Address:   client.DefaultHostPort,
			Namespace: client.DefaultNamespace,
		},
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}
