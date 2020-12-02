package export

import (
	"fmt"
	"io/ioutil"
	"main/src/config"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// prometheus.yml global
type PrometheusGlobal struct {
	ScrapeInterval     time.Duration `yaml:"scrape_interval"`
	EvaluationInterval time.Duration `yaml:"evaluation_interval"`
}

// prometheus.yml scrape_configs
type PrometheusScrapeConfig struct {
	JobName        string                    `yaml:"job_name"`
	ScrapeInterval time.Duration             `yaml:"scrape_interval"`
	MetricsPath    string                    `yaml:"metrics_path"`
	StaticConfigs  []*PrometheusStaticConfig `yaml:"static_configs"`
}

// prometheus.yml static_configs
type PrometheusStaticConfig struct {
	Targets []string `yaml:"targets"`
}

type PrometheusServiceExporter struct {
	Global        *PrometheusGlobal         `yaml:"global"`
	ScrapeConfigs []*PrometheusScrapeConfig `yaml:"scrape_configs"`
}

type PrometheusExporter struct {
	cse *PrometheusServiceExporter
}

func NewPrometheusExporter() *PrometheusExporter {
	return &PrometheusExporter{
		cse: &PrometheusServiceExporter{
			Global: &PrometheusGlobal{
				ScrapeInterval:     time.Second * 15,
				EvaluationInterval: time.Second * 15,
			},

			ScrapeConfigs: make([]*PrometheusScrapeConfig, 0),
		},
	}
}

func (ce *PrometheusExporter) WriteToFile(configs config.CombinedServices, path string) {
	mapServiceAddr := make(map[string]bool)

	for _, config := range configs {
		mapServiceAddr[fmt.Sprintf("%s:%s", config.InnerIP, config.WatcherPort)] = true
	}

	// generate watch_service
	for addr := range mapServiceAddr {
		scrapeConfig := &PrometheusScrapeConfig{
			JobName:        addr,
			ScrapeInterval: time.Second * 5,
			MetricsPath:    "/metrics",
			StaticConfigs:  make([]*PrometheusStaticConfig, 0),
		}

		scrapeConfig.StaticConfigs = append(scrapeConfig.StaticConfigs, &PrometheusStaticConfig{Targets: []string{addr}})
		ce.cse.ScrapeConfigs = append(ce.cse.ScrapeConfigs, scrapeConfig)
	}

	data, err := yaml.Marshal(ce.cse)
	if err != nil {
		log.Fatal().Err(err).Msg("yaml marshal failed")
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("write prometheus.yml failed")
	}

	log.Info().Str("path", path).Msg("write prometheus.yml successful")
}

// test
func (ce *PrometheusExporter) UnmarshalToStruct() {
	data, err := ioutil.ReadFile("../config/prometheus/prometheus.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("read prometheus.yml failed")
	}

	err = yaml.Unmarshal(data, ce.cse)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshal yaml failed")
	}

	log.Info().Interface("prometheus yaml", ce.cse).Msg("unmarshal success")
}
