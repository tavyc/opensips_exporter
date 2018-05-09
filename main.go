// Export OpenSIPS stats to Prometheus.

package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/tavyc/opensips_exporter/opensips_mi"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "opensips"

// OpensSIPS Prometheus exporter
type opensipsExporter struct {
	url string

	mu         sync.RWMutex
	commands   map[string]bool
	processes  [][]string
	profiles   map[string]bool
	lastUptime float64

	up                 *prometheus.Desc
	versionInfo        *prometheus.Desc
	processInfo        *prometheus.Desc
	profilesValuesInfo *prometheus.Desc
}

func (ose *opensipsExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- ose.up
	ch <- ose.versionInfo
	ch <- ose.processInfo
	ch <- ose.profilesValuesInfo

	for _, stats := range opensipsStats {
		for _, stat := range stats {
			ch <- stat.desc
		}
	}
}

func (ose *opensipsExporter) Collect(ch chan<- prometheus.Metric) {
	up := 0

	defer (func() {
		ch <- prometheus.MustNewConstMetric(ose.up, prometheus.GaugeValue, float64(up))
	})()

	conn, err := opensips_mi.NewMIJsonClient(ose.url, opensips_mi.MIJsonConfig{})
	if err != nil {
		log.Print("error connecting to OpensSIPS: ", err)
		return
	}
	defer conn.Close()

	if err = ose.collectVersionInfo(conn, ch); err != nil {
		return
	}

	var uptime float64
	up = 1

	ose.mu.RLock()
	hasCommands := len(ose.commands) > 0
	hasProcesses := len(ose.processes) > 0
	hasProfiles := len(ose.profiles) > 0
	ose.mu.RUnlock()

	if !hasCommands {
		ose.fetchCommands(conn)
	}

	ose.mu.RLock()
	hasStatisticsCommand := ose.commands["get_statistics"]
	hasProfilesCommand := ose.commands["list_all_profiles"]
	ose.mu.RUnlock()

	ose.collectProcessInfo(conn, ch, !hasProcesses)
	if hasStatisticsCommand {
		uptime = ose.collectStats(conn, ch)
	}
	if hasProfilesCommand {
		ose.collectDialogProfiles(conn, ch, !hasProfiles)
	}

	// Invalidate our caches when the monitored target restarts
	ose.mu.RLock()
	restart := uptime < ose.lastUptime
	ose.mu.RUnlock()

	if restart {
		ose.mu.Lock()
		ose.commands = make(map[string]bool)
		ose.processes = nil
		ose.profiles = make(map[string]bool)
		ose.lastUptime = uptime
		ose.mu.Unlock()
	}
}

var versionRegexp = regexp.MustCompile(`(\S+)\s+\((\S+)\s+\((\S+)/(\S+)\)\)`)

func (ose *opensipsExporter) collectVersionInfo(conn opensips_mi.Client, ch chan<- prometheus.Metric) error {
	resp, err := conn.Command("version")
	if err != nil {
		return err
	}
	m := versionRegexp.FindStringSubmatch(resp.ChildValues["Server"])
	if m != nil {
		ch <- prometheus.MustNewConstMetric(ose.versionInfo, prometheus.GaugeValue, 1, m[1:]...)
	}
	return nil
}

func (ose *opensipsExporter) fetchCommands(conn opensips_mi.Client) {
	resp, err := conn.Command("which")
	if err != nil {
		return
	}
	cmds := make(map[string]bool, len(resp.Children))
	for _, node := range resp.Children {
		cmds[node.Value] = true
	}

	ose.mu.Lock()
	ose.commands = cmds
	ose.mu.Unlock()
}

func (ose *opensipsExporter) collectProcessInfo(conn opensips_mi.Client, ch chan<- prometheus.Metric, update bool) {
	var processes [][]string

	if update {
		resp, err := conn.Command("ps")
		if err != nil {
			return
		}
		processes = make([][]string, 0, len(resp.Children))
		for _, node := range resp.Children {
			processes = append(processes, []string{node.Attrs["ID"], strings.TrimSpace(node.Attrs["Type"])})
		}

		ose.mu.Lock()
		ose.processes = processes
		ose.mu.Unlock()
	}

	ose.mu.RLock()
	defer ose.mu.RUnlock()

	for _, proc := range ose.processes {
		ch <- prometheus.MustNewConstMetric(ose.processInfo, prometheus.GaugeValue, 1, proc...)
	}
}

func (ose *opensipsExporter) collectStats(conn opensips_mi.Client, ch chan<- prometheus.Metric) (uptime float64) {
	resp, err := conn.Command("get_statistics", "all")
	if err != nil {
		return
	}
	for statName, statValue := range resp.ChildValues {
		parts := strings.SplitN(statName, ":", 2)
		if len(parts) != 2 {
			continue
		}
		subsys := parts[0]
		metric := strings.Replace(parts[1], " ", "_", -1)
		value, err := strconv.ParseFloat(statValue, 64)
		if err != nil {
			continue
		}

		if statName == "core:timestamp" {
			uptime = value
		}

		stats, exists := opensipsStats[subsys]
		if !exists {
			continue
		}

		for _, stat := range stats {
			if stat.regexp != nil {
				mm := stat.regexp.FindStringSubmatch(metric)
				if mm != nil {
					ch <- prometheus.MustNewConstMetric(stat.desc, stat.value, value, mm[1:]...)
					break
				}
			} else if metric == stat.stat {
				ch <- prometheus.MustNewConstMetric(stat.desc, stat.value, value)
				break
			}
		}
	}

	return
}

var profileValuesRegexp = regexp.MustCompile(`(?:^|,)([a-z0-9_]+)=([^,]*)`)

func (ose *opensipsExporter) collectDialogProfiles(conn opensips_mi.Client, ch chan<- prometheus.Metric, update bool) {
	var profiles map[string]bool

	if update {
		resp, err := conn.Command("list_all_profiles")
		if err != nil {
			return
		}

		profiles = make(map[string]bool, len(resp.ChildValues))
		for profile, hasValues := range resp.ChildValues {
			profiles[profile] = hasValues != "0"
		}
		ose.mu.Lock()
		ose.profiles = profiles
		ose.mu.Unlock()
	}

	ose.mu.RLock()
	defer ose.mu.RUnlock()

	for profile, hasValues := range ose.profiles {
		if !hasValues {
			continue
		}

		getResp, err := conn.Command("profile_get_values", profile)
		if err != nil {
			continue
		}

		for _, node := range getResp.Children {
			count, err := strconv.ParseFloat(node.Attrs["count"], 64)
			if err != nil {
				continue
			}

			// Parse dialog value as "name=value," pairs and export the pairs as labels
			matches := profileValuesRegexp.FindAllStringSubmatch(node.Value, -1)
			if matches != nil {
				labelNames := []string{"profile"}
				labels := []string{profile}
				for _, match := range matches {
					labelNames = append(labelNames, match[1])
					labels = append(labels, match[2])
				}
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						prometheus.BuildFQName(namespace, "dialog", "profiles_with_values_count"),
						"Dialog profiles with counts",
						labelNames,
						nil,
					),
					prometheus.GaugeValue,
					count,
					labels...,
				)
			} else {
				// Export just the profile and value labels
				ch <- prometheus.MustNewConstMetric(ose.profilesValuesInfo, prometheus.GaugeValue, count,
					profile, node.Value)
			}
		}
	}
}

func newOpensipsExporter(url string) *opensipsExporter {
	return &opensipsExporter{
		url: url,

		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"1 if OpenSIPS is running",
			nil,
			nil,
		),
		versionInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "version_info"),
			"Version information (always 1)",
			[]string{"server", "version", "arch", "os"},
			nil,
		),
		processInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "process_info"),
			"Process information (always 1)",
			[]string{"id", "type"},
			nil,
		),
		profilesValuesInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "dialog", "profiles_with_values_count"),
			"Dialog profiles with counts",
			[]string{"profile", "value"},
			nil,
		),
	}
}

var (
	url = flag.String("opensips.url", "http://127.0.0.1:8062/json",
		"The HTTP address to connect to OpenSIPS mi_json")
	listenAddr = flag.String("web.listen-address", ":9441",
		"The address to listen on for HTTP requests.")
)

func main() {
	flag.Parse()

	prometheus.MustRegister(newOpensipsExporter(*url))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
