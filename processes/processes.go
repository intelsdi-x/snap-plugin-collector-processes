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
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/serror"

	"github.com/intelsdi-x/snap-plugin-utilities/str"
)

const (
	//PLUGIN name
	PLUGIN = "processes"
	//VENDOR name
	VENDOR = "intel"
	// FS is proc filesystem
	FS = "procfs"
	//VERSION of plugin
	VERSION = 3
)

var (
	metricNames = []string{
		"ps_vm",
		"ps_rss",
		"ps_data",
		"ps_code",
		"ps_stacksize",
		"ps_cputime_user",
		"ps_cputime_system",
		"ps_pagefaults_min",
		"ps_pagefaults_maj",
		"ps_disk_ops_syscr",
		"ps_disk_ops_syscw",
		"ps_disk_octets_rchar",
		"ps_disk_octets_wchar",
		"ps_count",
	}
)

// New returns instance of processes plugin
func New() *procPlugin {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	return &procPlugin{host: host, mc: &procStatsCollector{}}
}

// GetMetricTypes returns list of available metrics
func (procPlg *procPlugin) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	metricTypes := []plugin.MetricType{}
	// build metric types from process metric names
	for _, metricName := range metricNames {
		metricType := plugin.MetricType{
			Namespace_: core.NewNamespace(VENDOR, FS, PLUGIN).
				AddDynamicElement("process_name", "name of the running process").
				AddStaticElements(metricName),
			Config_: cfg.ConfigDataNode,
		}
		metricTypes = append(metricTypes, metricType)
	}
	// build metric types from process states
	for _, state := range States.Values() {
		metricType := plugin.MetricType{
			Namespace_: core.NewNamespace(VENDOR, FS, PLUGIN, state),
			Config_:    cfg.ConfigDataNode,
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
func (procPlg *procPlugin) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	metrics := []plugin.MetricType{}
	stateCount := map[string]int{}

	// init stateCount map with keys from States
	for _, state := range States.Values() {
		stateCount[state] = 0
	}
	// get all proc stats
	stats, err := procPlg.mc.GetStats()
	if err != nil {
		return nil, serror.New(err)
	}
	// calculate number of proces in each state
	for _, instances := range stats {
		for _, instance := range instances {
			stateName := States[instance.State]
			stateCount[stateName]++
		}
	}
	// calculate metrics
	for _, metricType := range metricTypes {
		ns := metricType.Namespace()
		if len(ns) < 4 {
			return nil, serror.New(fmt.Errorf("Unknown namespace length. Expecting at least 4, is %d", len(ns)))
		}
		// get incoming tags and add hostname
		tags := metricType.Tags()
		if tags == nil {
			tags = map[string]string{}
		}
		tags["hostname"] = procPlg.host

		isDynamic, _ := ns.IsDynamic()
		if isDynamic {
			// ns[3] is wildcard for all processes
			for procName, instances := range stats {
				procMetrics := setProcMetrics(instances)
				for procMet, val := range procMetrics {
					if procMet == ns[4].Value {
						metric := plugin.MetricType{
							Namespace_: core.NewNamespace(VENDOR, FS, PLUGIN, procName, procMet),
							Data_:      val,
							Timestamp_: time.Now(),
							Tags_:      tags,
						}
						metrics = append(metrics, metric)
					}
				}
			}
		} else if str.Contains(States.Values(), ns[3].Value) {
			// ns[3] contains process state
			state := ns[3].Value
			if val, ok := stateCount[state]; ok {
				metric := plugin.MetricType{
					Namespace_: ns,
					Data_:      val,
					Timestamp_: time.Now(),
					Tags_:      tags,
				}
				metrics = append(metrics, metric)
			}
		} else {
			// ns[3] contains process name
			if len(ns) != 5 {
				return nil, serror.New(fmt.Errorf("Unknown namespace length. Expecting at 5, is %d", len(ns)))
			}
			procName := ns[3].Value
			metricName := ns[4].Value
			instances, found := stats[procName]
			if !found {
				return nil, serror.New(fmt.Errorf("Process name {%s} not found!", procName))
			}
			procMetrics := setProcMetrics(instances)
			for procMet, val := range procMetrics {
				if metricName == procMet {
					metric := plugin.MetricType{
						Namespace_: core.NewNamespace(VENDOR, FS, PLUGIN, procName, procMet),
						Data_:      val,
						Timestamp_: time.Now(),
						Tags_:      tags,
					}
					metrics = append(metrics, metric)
					break
				}
			}
		}
	}

	return metrics, nil
}

func setProcMetrics(instances []Proc) map[string]uint64 {
	procMetrics := map[string]uint64{}
	for _, metricName := range metricNames {
		procMetrics[metricName] = 0
	}
	procMetrics["ps_count"] = uint64(len(instances))

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

	return procMetrics
}

// procPlugin holds host name and reference to metricCollector which has method of GetStats()
type procPlugin struct {
	host string
	mc   metricCollector
}
