// +build linux

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

package processes

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
)

const (
	//PLUGIN name
	PLUGIN = "processes"
	//VENDOR name
	VENDOR = "intel"
	// FS is proc filesystem
	FS = "procfs"
	//VERSION of plugin
	VERSION = 1
)

// procPlugin holds host name and reference to metricCollector which has method of GetStats()
type procPlugin struct {
	host string
	mc   metricCollector
}

// New returns instance of processes plugin
func New() *procPlugin {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	return &procPlugin{host: host, mc: &procStatsCollector{}}
}

// GetMetricTypes returns list of available metrics
func (procPlg *procPlugin) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	metricTypes := []plugin.PluginMetricType{
		plugin.PluginMetricType{
			Namespace_: []string{VENDOR, FS, PLUGIN, "*"},
		},
	}
	for _, state := range States {
		metricType := plugin.PluginMetricType{
			Namespace_: []string{VENDOR, FS, PLUGIN, state},
		}
		metricTypes = append(metricTypes, metricType)
	}
	return metricTypes, nil
}

// GetConfigPolicy returns config policy
func (procPlg *procPlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

// CollectMetrics retrieves values for given metrics types
func (procPlg *procPlugin) CollectMetrics(metricTypes []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := []plugin.PluginMetricType{}
	stateCount := map[string]int{}

	// init stateCount map with keys from States
	for _, state := range States {
		stateCount[state] = 0
	}

	stats, err := procPlg.mc.GetStats()

	if err != nil {
		return nil, err
	}
	for _, instances := range stats {
		for _, instance := range instances {
			stateName := States[instance.State]
			stateCount[stateName]++
		}
	}

	for _, metricType := range metricTypes {
		ns := metricType.Namespace()
		if len(ns) != 4 {
			return nil, fmt.Errorf("Unknown namespace length. Expecting 4, is %d", len(ns))
		}
		state := ns[3]
		if state == "*" {
			for procName, instances := range stats {
				procMetrics := map[string]uint64{
					"ps_vm":                0,
					"ps_rss":               0,
					"ps_data":              0,
					"ps_code":              0,
					"ps_stacksize":         0,
					"ps_cputime_user":      0,
					"ps_cputime_system":    0,
					"ps_pagefaults_min":    0,
					"ps_pagefaults_maj":    0,
					"ps_disk_ops_syscr":    0,
					"ps_disk_ops_syscw":    0,
					"ps_disk_octets_rchar": 0,
					"ps_disk_octets_wchar": 0,
					"ps_count":             uint64(len(instances)),
				}

				for _, instance := range instances {
					vm, _ := strconv.ParseUint(string(instance.Stat[22]), 10, 64)
					procMetrics["ps_vm"] += vm

					rss, _ := strconv.ParseUint(string(instance.Stat[23]), 10, 64)
					procMetrics["ps_rss"] += rss

					procMetrics["ps_data"] += instance.VmData
					procMetrics["ps_code"] += instance.VmCode

					stack1, _ := strconv.ParseUint(string(instance.Stat[27]), 10, 64)
					stack2, _ := strconv.ParseUint(string(instance.Stat[28]), 10, 64)

					// to avoid overload
					if stack1 > stack2 {
						procMetrics["ps_stacksize"] += stack1 - stack2
					} else {
						procMetrics["ps_stacksize"] += stack2 - stack1
					}

					utime, _ := strconv.ParseUint(string(instance.Stat[13]), 10, 64)
					procMetrics["ps_cputime_user"] += utime * 10000

					stime, _ := strconv.ParseUint(string(instance.Stat[14]), 10, 64)
					procMetrics["ps_cputime_system"] += stime * 10000

					minflt, _ := strconv.ParseUint(string(instance.Stat[9]), 10, 64)
					procMetrics["ps_pagefaults_min"] += minflt * 10000

					majflt, _ := strconv.ParseUint(string(instance.Stat[11]), 10, 64)
					procMetrics["ps_pagefaults_maj"] += majflt * 10000

					procMetrics["ps_disk_octets_rchar"] += instance.Io["rchar"]
					procMetrics["ps_disk_octets_wchar"] += instance.Io["wchar"]
					procMetrics["ps_disk_ops_syscr"] += instance.Io["syscr"]
					procMetrics["ps_disk_ops_syscw"] += instance.Io["syscw"]
				}

				for procMet, val := range procMetrics {
					metric := plugin.PluginMetricType{
						Namespace_: []string{VENDOR, FS, PLUGIN, procName, procMet},
						Data_:      val,
						Timestamp_: time.Now(),
						Source_:    procPlg.host,
					}
					metrics = append(metrics, metric)
				}
				// TODO - if ( report_ctx_switch ) ps_read_tasks_status(pid, ps)
			}
		} else {
			if val, ok := stateCount[state]; ok {
				metric := plugin.PluginMetricType{
					Namespace_: ns,
					Data_:      val,
					Timestamp_: time.Now(),
					Source_:    procPlg.host,
				}
				metrics = append(metrics, metric)
			}
		}

	}
	return metrics, nil
}
