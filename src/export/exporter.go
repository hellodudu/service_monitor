package export

import (
	"encoding/json"
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
	MetricsPath    string                    `yaml:"metrics_path,omitempty"`
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
	cse  *PrometheusServiceExporter
	host []*HostJsonConfig
}

type HostJsonLabels struct {
	Group    string `json:"group"`
	App      string `json:"app"`
	Hostname string `json:"hostname"`
}

type HostJsonConfig struct {
	Targets []string        `json:"targets"`
	Labels  *HostJsonLabels `json:"labels"`
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
		host: make([]*HostJsonConfig, 0),
	}
}

func (ce *PrometheusExporter) WriteToFile(configs config.CombinedServices, path string) {
	mapServiceAddr := make(map[string]string)
	mapServiceNode := make(map[string]string)

	for _, config := range configs {
		mapServiceAddr[fmt.Sprintf("%s:%s", config.InnerIP, config.WatcherPort)] = config.ProcessName
		mapServiceNode[fmt.Sprintf("%s:%s", config.InnerIP, config.NodePort)] = config.ProcessName
	}

	// generate watch_service
	// for addr, proName := range mapServiceAddr {
	// 	scrapeConfig := &PrometheusScrapeConfig{
	// 		JobName:        proName,
	// 		ScrapeInterval: time.Second * 5,
	// 		MetricsPath:    "/metrics",
	// 		StaticConfigs:  make([]*PrometheusStaticConfig, 0),
	// 	}

	// 	scrapeConfig.StaticConfigs = append(scrapeConfig.StaticConfigs, &PrometheusStaticConfig{Targets: []string{addr}})
	// 	ce.cse.ScrapeConfigs = append(ce.cse.ScrapeConfigs, scrapeConfig)
	// }

	// // generate node_exporter
	// for addr, proName := range mapServiceNode {
	// 	scrapeConfig := &PrometheusScrapeConfig{
	// 		JobName:        fmt.Sprintf("node_%s", proName),
	// 		ScrapeInterval: time.Second * 5,
	// 		MetricsPath:    "/metrics",
	// 		StaticConfigs:  make([]*PrometheusStaticConfig, 0),
	// 	}

	// 	scrapeConfig.StaticConfigs = append(scrapeConfig.StaticConfigs, &PrometheusStaticConfig{Targets: []string{addr}})
	// 	ce.cse.ScrapeConfigs = append(ce.cse.ScrapeConfigs, scrapeConfig)
	// }

	// data, err := yaml.Marshal(ce.cse)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("yaml marshal failed")
	// }

	for addr, proName := range mapServiceAddr {
		hostConfig := &HostJsonConfig{
			Targets: []string{addr},
			Labels: &HostJsonLabels{
				Group:    "ET服务",
				App:      "c#",
				Hostname: fmt.Sprintf("%sET服务", proName),
			},
		}

		ce.host = append(ce.host, hostConfig)
	}

	for addr, proName := range mapServiceNode {
		hostConfig := &HostJsonConfig{
			Targets: []string{addr},
			Labels: &HostJsonLabels{
				Group:    "物理机节点",
				App:      "node",
				Hostname: fmt.Sprintf("%s节点", proName),
			},
		}

		ce.host = append(ce.host, hostConfig)
	}

	data, err := json.Marshal(ce.host)
	if err != nil {
		log.Fatal().Err(err).Msg("host config marshal failed")
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("write host.config failed")
	}

	log.Info().Str("path", path).Msg("write host.config successful")
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
