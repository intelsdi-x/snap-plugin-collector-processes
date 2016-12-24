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
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"strings"
)

const (
	//plugin name
	PluginName = "processes"
	//plugin vendor
	PluginVendor = "intel"
	// fs is proc filesystem
	fs = "procfs"
	//version of plugin
	PluginVersion = 8
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
		"ps_cmd_line": label{
			description: "Process command line with full path and args",
			unit:        "",
		},
		"ps_cmd": label{
			description: "Process command line with full path",
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

// GetMetricTypes returns list of available metrics
func (procPlg *procPlugin) GetMetricTypes(_ plugin.Config) ([]plugin.Metric, error) {
	metricTypes := []plugin.Metric{}
	// build metric types from process metric names
	for metricName, label := range metricNames {
		metricType := plugin.Metric{
			Namespace: plugin.NewNamespace(PluginVendor, fs, PluginName, "pid").
				AddDynamicElement("process_id", "pid of the running process").
				AddStaticElement(metricName),
			Description: label.description,
			Unit:        label.unit,
		}
		metricTypes = append(metricTypes, metricType)
		metricType2 := plugin.Metric{
			Namespace: plugin.NewNamespace(PluginVendor, fs, PluginName, "name").
				AddDynamicElement("name", "name of the running process").
				AddDynamicElement("process_id", "pid of the running process").
				AddStaticElement(metricName),
			Description: label.description,
			Unit:        label.unit,
		}
		metricTypes = append(metricTypes, metricType2)

	}
	// build metric types from process states
	for _, state := range States.Values() {
		metricType := plugin.Metric{
			Namespace:   plugin.NewNamespace(PluginVendor, fs, PluginName, "stats", state),
			Description: fmt.Sprintf("Number of processes in %s state", state),
			Unit:        "",
		}
		metricTypes = append(metricTypes, metricType)
	}
	return metricTypes, nil
}

// GetConfigPolicy returns config policy
func (procPlg *procPlugin) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	configKey := []string{PluginVendor, fs, PluginName}
	policy.AddNewStringRule(configKey, "proc_path", false, plugin.SetDefaultString("/proc"))
	policy.AddNewBoolRule(configKey, "include_system_processes", false, plugin.SetDefaultBool(false))
	return *policy, nil
}

// CollectMetrics retrieves values for given metrics types
func (procPlg *procPlugin) CollectMetrics(metric []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	stateCount := map[string]int{}

	//procPath, err := config.GetConfigItem(metricTypes[0], "proc_path")
	procPath, err := metric[0].Config.GetString("proc_path")
	if err != nil {
		return nil, err
	}
	isp, err := metric[0].Config.GetBool("include_system_processes")
	if err != nil {
		return nil, err
	}

	// init stateCount map with keys from States

	for _, state := range States.Values() {
		stateCount[state] = 0
	}
	timestamp := time.Now()
	// get all proc stats
	stats, err := procPlg.mc.GetStats(procPath, isp)
	if err != nil {
		return nil, err
	}
	// calculate number of processes in each state

	for _, proc := range stats {
		stateName := States[proc.State]
		stateCount[stateName]++
	}

	// calculate metrics
	for _, metric := range metric {
		ns := metric.Namespace
		if len(ns) < 4 {
			return nil, errors.New("Unknown namespace length. Expecting at least 4, is " + strconv.Itoa(len(ns)))
		}

		isDynamic, _ := ns.IsDynamic()
		if isDynamic {
			for _, proc := range stats {
				procMetrics := setProcMetrics(proc)
				for procMet, val := range procMetrics {
					switch ns[3].Value {
					//[0][1][2] = /intel/procfs/processes, [3] = "pid", [4] = <pid>, [5] = metric name
					case "pid":
						if procMet == ns[5].Value {
							//nuns := plugin.Namespace(append([]plugin.NamespaceElement{}, ns...))
							nsClone := make([]plugin.NamespaceElement, len(ns))
							copy(nsClone, ns)
							nsClone[4].Value = strconv.Itoa(proc.Pid)
							newMetric := plugin.Metric{
								Namespace:   nsClone,
								Data:        val,
								Timestamp:   timestamp,
								Unit:        metricNames[procMet].unit,
								Description: metricNames[procMet].description,
							}
							metrics = append(metrics, newMetric)

						}
					case "name":
						if procMet == ns[6].Value {
							//nuns := plugin.Namespace(append([]plugin.NamespaceElement{}, ns...))
							nsClone := make([]plugin.NamespaceElement, len(ns))
							copy(nsClone, ns)
							nsClone[4].Value = removeUnwantedChars(proc.Cmd)
							nsClone[5].Value = strconv.Itoa(proc.Pid)
							newMetric := plugin.Metric{
								Namespace:   nsClone,
								Data:        val,
								Timestamp:   timestamp,
								Unit:        metricNames[procMet].unit,
								Description: metricNames[procMet].description,
							}
							metrics = append(metrics, newMetric)

						}
					}
				}
			}
		} else if contains(States.Values(), ns[4].Value) {
			// ns[3] contains process state
			state := ns[4].Value
			if val, ok := stateCount[state]; ok {
				metric := plugin.Metric{
					Namespace:   ns,
					Data:        val,
					Timestamp:   timestamp,
					Description: fmt.Sprintf("Number of processes in %s state", state),
					Unit:        "",
				}
				metrics = append(metrics, metric)
			}
		}
	}

	return metrics, nil
}

func setProcMetrics(proc Proc) map[string]interface{} {

	procMetrics := map[string]interface{}{}

	for metricName, _ := range metricNames {
		procMetrics[metricName] = 0
	}
	vm, _ := strconv.ParseUint(string(proc.Stat[22]), 10, 64)
	procMetrics["ps_vm"] = vm

	rss, _ := strconv.ParseUint(string(proc.Stat[23]), 10, 64)
	procMetrics["ps_rss"] = rss

	procMetrics["ps_data"] = proc.VmData
	procMetrics["ps_code"] = proc.VmCode

	stack1, _ := strconv.ParseUint(string(proc.Stat[27]), 10, 64)
	stack2, _ := strconv.ParseUint(string(proc.Stat[28]), 10, 64)

	// to avoid overload
	if stack1 > stack2 {
		procMetrics["ps_stacksize"] = stack1 - stack2
	} else {
		procMetrics["ps_stacksize"] = stack2 - stack1
	}

	utime, _ := strconv.ParseUint(string(proc.Stat[13]), 10, 64)
	procMetrics["ps_cputime_user"] = utime

	stime, _ := strconv.ParseUint(string(proc.Stat[14]), 10, 64)
	procMetrics["ps_cputime_system"] = stime

	minflt, _ := strconv.ParseUint(string(proc.Stat[9]), 10, 64)
	procMetrics["ps_pagefaults_min"] = minflt

	majflt, _ := strconv.ParseUint(string(proc.Stat[11]), 10, 64)
	procMetrics["ps_pagefaults_maj"] = majflt

	procMetrics["ps_disk_octets_rchar"] = proc.Io["rchar"]
	procMetrics["ps_disk_octets_wchar"] = proc.Io["wchar"]
	procMetrics["ps_disk_ops_syscr"] = proc.Io["syscr"]
	procMetrics["ps_disk_ops_syscw"] = proc.Io["syscw"]
	procMetrics["ps_cmd_line"] = proc.CmdLine
	procMetrics["ps_cmd"] = proc.Cmd
	//	}

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

func contains(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}

func removeUnwantedChars(str string) string {
	unwanteds := []unwanted{
		{"[", ""},
		{"]", ""},
		{"(", ""},
		{")", ""},
		{"/", "."},
		{"\\", ""},
	}
	for _, unw := range unwanteds {
		str = strings.Replace(str, unw.char, unw.repl, -1)
	}
	return str
}
