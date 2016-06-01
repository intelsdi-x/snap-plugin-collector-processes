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

	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap-plugin-utilities/str"
)

const (
	//plugin name
	pluginName = "processes"
	//plugin vendor
	pluginVendor = "intel"
	// fs is proc filesystem
	fs = "procfs"
	//version of plugin
	version = 4
)

var (
	metricNames = map[string]label{
		"ps_vm": label{
			description: "Virtual memory size in bytes",
			unit:        "B",
		},
		"ps_rss": label{
			description: "Resident Set Size: number of pages the process has in real memory",
			unit:        "",
		},
		"ps_data": label{
			description: "Size of data segments",
			unit:        "B",
		},
		"ps_code": label{
			description: "Size of text segment",
			unit:        "B",
		},
		"ps_stacksize": label{
			description: "Stack size",
			unit:        "B",
		},
		"ps_cputime_user": label{
			description: "Amount of time that this process has been scheduled in user mode",
			unit:        "Jiff",
		},
		"ps_cputime_system": label{
			description: "Amount of time that this process has been scheduled in kernel mode",
			unit:        "Jiff",
		},
		"ps_pagefaults_min": label{
			description: "The number of minor faults the process has made",
			unit:        "",
		},
		"ps_pagefaults_maj": label{
			description: "The number of major faults the process has made",
			unit:        "",
		},
		"ps_disk_ops_syscr": label{
			description: "Attempt to count the number of read I/O operations",
			unit:        "",
		},
		"ps_disk_ops_syscw": label{
			description: "Attempt to count the number of write I/O operations",
			unit:        "",
		},
		"ps_disk_octets_rchar": label{
			description: "The number of bytes which this task has caused to be read from storage",
			unit:        "B",
		},
		"ps_disk_octets_wchar": label{
			description: "The number of bytes which this task has caused, or shall cause to be written to disk",
			unit:        "B",
		},
		"ps_count": label{
			description: "Number of process instances",
			unit:        "",
		},
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

// Meta returns plugin meta data
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		pluginVendor,
		version,
		plugin.CollectorPluginType,
		[]string{},
		[]string{plugin.SnapGOBContentType},
		plugin.ConcurrencyCount(1),
	)
}

// GetMetricTypes returns list of available metrics
func (procPlg *procPlugin) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	metricTypes := []plugin.MetricType{}
	// build metric types from process metric names
	for metricName, label := range metricNames {
		metricType := plugin.MetricType{
			Namespace_: core.NewNamespace(pluginVendor, fs, pluginName).
				AddDynamicElement("process_name", "name of the running process").
				AddStaticElements(metricName),
			Config_:      cfg.ConfigDataNode,
			Description_: label.description,
			Unit_:        label.unit,
		}
		metricTypes = append(metricTypes, metricType)
	}
	// build metric types from process states
	for _, state := range States.Values() {
		metricType := plugin.MetricType{
			Namespace_:   core.NewNamespace(pluginVendor, fs, pluginName, state),
			Config_:      cfg.ConfigDataNode,
			Description_: fmt.Sprintf("Number of processes in %s state", state),
			Unit_:        "",
		}
		metricTypes = append(metricTypes, metricType)
	}
	return metricTypes, nil
}

// GetConfigPolicy returns config policy
func (procPlg *procPlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	rule, _ := cpolicy.NewStringRule("procfs_path", false, "/proc")
	node := cpolicy.NewPolicyNode()
	node.Add(rule)
	cp.Add([]string{pluginVendor, fs, pluginName}, node)
	return cp, nil
}

// CollectMetrics retrieves values for given metrics types
func (procPlg *procPlugin) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	metrics := []plugin.MetricType{}
	stateCount := map[string]int{}

	procPath, err := config.GetConfigItem(metricTypes[0], "procfs_path")
	if err != nil {
		return nil, err
	}

	// init stateCount map with keys from States
	for _, state := range States.Values() {
		stateCount[state] = 0
	}
	// get all proc stats
	stats, err := procPlg.mc.GetStats(procPath.(string))
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

		isDynamic, _ := ns.IsDynamic()
		if isDynamic {
			// ns[3] is wildcard for all processes
			for procName, instances := range stats {
				procMetrics := setProcMetrics(instances)
				for procMet, val := range procMetrics {
					if procMet == ns[4].Value {
						// change dynamic namespace element value (= "*") to current process name
						// whole namespace stays dynamic (ns[3].Name != "")
						ns[3].Value = procName
						metric := plugin.MetricType{
							Namespace_:   ns,
							Data_:        val,
							Timestamp_:   time.Now(),
							Unit_:        metricNames[procMet].unit,
							Description_: metricNames[procMet].description,
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
					Namespace_:   ns,
					Data_:        val,
					Timestamp_:   time.Now(),
					Description_: fmt.Sprintf("Number of processes in %s state", state),
					Unit_:        "",
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
						Namespace_:   core.NewNamespace(pluginVendor, fs, pluginName, procName, procMet),
						Data_:        val,
						Timestamp_:   time.Now(),
						Description_: metricNames[procMet].description,
						Unit_:        metricNames[procMet].unit,
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
	for metricName, _ := range metricNames {
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
		procMetrics["ps_cputime_user"] += utime

		stime, _ := strconv.ParseUint(string(instance.Stat[14]), 10, 64)
		procMetrics["ps_cputime_system"] += stime

		minflt, _ := strconv.ParseUint(string(instance.Stat[9]), 10, 64)
		procMetrics["ps_pagefaults_min"] += minflt

		majflt, _ := strconv.ParseUint(string(instance.Stat[11]), 10, 64)
		procMetrics["ps_pagefaults_maj"] += majflt

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

type label struct {
	description string
	unit        string
}
