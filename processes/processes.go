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
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	// PluginName is name of plugin
	PluginName = "processes"
	// PluginVersion is version of plugin
	PluginVersion = 8

	//plugin vendor
	pluginVendor = "intel"
	// fs is proc filesystem
	fs = "procfs"

	// Namespace offsets
	nsCategory  = 3
	nsProcName  = 4 // /intel/procfs/processes/process/->ProcName<-
	nsStateName = 4 // /intel/procfs/processes/states/->StateName<-
	nsPid       = 5 // /intel/procfs/processes/process/ProcName/->Pid<-
	nsPsCount   = 5 // /intel/procfs/processes/process/ProcName/->ps_count<-
	nsPidMetric = 6 // /intel/procfs/processes/process/ProcName/Pid/->metric<-
)

var (
	metricNames = map[string]label{
		"ps_vm": label{
			category:    "pid",
			description: "Virtual memory size in bytes",
			unit:        "B",
		},
		"ps_rss": label{
			category:    "pid",
			description: "Resident Set Size: number of pages the process has in real memory",
		},
		"ps_data": label{
			category:    "pid",
			description: "Size of data segments",
			unit:        "B",
		},
		"ps_code": label{
			category:    "pid",
			description: "Size of text segment",
			unit:        "B",
		},
		"ps_stacksize": label{
			category:    "pid",
			description: "Stack size",
			unit:        "B",
		},
		"ps_cputime_user": label{
			category:    "pid",
			description: "Amount of time that this process has been scheduled in user mode",
			unit:        "Jiff",
		},
		"ps_cputime_system": label{
			category:    "pid",
			description: "Amount of time that this process has been scheduled in kernel mode",
			unit:        "Jiff",
		},
		"ps_pagefaults_min": label{
			category:    "pid",
			description: "The number of minor faults the process has made",
		},
		"ps_pagefaults_maj": label{
			category:    "pid",
			description: "The number of major faults the process has made",
		},
		"ps_disk_ops_syscr": label{
			category:    "pid",
			description: "Attempt to count the number of read I/O operations",
		},
		"ps_disk_ops_syscw": label{
			category:    "pid",
			description: "Attempt to count the number of write I/O operations",
		},
		"ps_disk_octets_rchar": label{
			category:    "pid",
			description: "The number of bytes which this task has caused to be read from storage",
			unit:        "B",
		},
		"ps_disk_octets_wchar": label{
			category:    "pid",
			description: "The number of bytes which this task has caused, or shall cause to be written to disk",
			unit:        "B",
		},
		"ps_cmdline": label{
			category:    "pid",
			description: "Process command line with arguments",
		},

		"ps_count": label{
			category:    "process",
			description: "Number of process instances",
		},

		"running": label{
			category:    "state",
			description: "Number of processes with 'running' status",
		},
		"sleeping": label{
			category:    "state",
			description: "Number of processes with 'sleeping' status",
		},
		"waiting": label{
			category:    "state",
			description: "Number of processes with 'waiting' status",
		},
		"zombie": label{
			category:    "state",
			description: "Number of processes with 'zombie' status",
		},
		"stopped": label{
			category:    "state",
			description: "Number of processes with 'stopped' status",
		},
		"tracing": label{
			category:    "state",
			description: "Number of processes with 'tracing' status",
		},
		"dead": label{
			category:    "state",
			description: "Number of processes with 'dead' status",
		},
		"wakekill": label{
			category:    "state",
			description: "Number of processes with 'wakekill' status",
		},
		"waking": label{
			category:    "state",
			description: "Number of processes with 'waking' status",
		},
		"parked": label{
			category:    "state",
			description: "Number of processes with 'parked' status",
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
func Meta() []plugin.MetaOpt {
	return []plugin.MetaOpt{
		plugin.ConcurrencyCount(1),
	}
}

// GetMetricTypes returns list of available metrics
func (procPlg *procPlugin) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	metricTypes := []plugin.Metric{}
	// build metric types from process metric names
	for metricName, label := range metricNames {
		switch label.category {
		case "pid":
			metricTypes = append(metricTypes, plugin.Metric{
				Namespace: plugin.NewNamespace(pluginVendor, fs, PluginName, "process").
					AddDynamicElement("process_name", "name of the process").
					AddDynamicElement("process_pid", "identifier of the process").
					AddStaticElements(metricName),
				Config:      cfg,
				Description: label.description,
				Unit:        label.unit,
			})

			// Aggregated metrics
			if metricName != "ps_cmdline" {
				metricTypes = append(metricTypes,
					plugin.Metric{
						Namespace: plugin.NewNamespace(pluginVendor, fs, PluginName, "process").
							AddDynamicElement("process_name", "name of the process").
							AddStaticElement("all").
							AddStaticElements(metricName),
						Config:      cfg,
						Description: label.description,
						Unit:        label.unit,
					})
			}
		case "state":
			metricTypes = append(metricTypes, plugin.Metric{
				Namespace:   plugin.NewNamespace(pluginVendor, fs, PluginName, "state", metricName),
				Config:      cfg,
				Description: label.description,
				Unit:        label.unit,
			})
		case "process":
			metricTypes = append(metricTypes,
				plugin.Metric{
					Namespace: plugin.NewNamespace(pluginVendor, fs, PluginName, "process").
						AddDynamicElement("process_name", "name of the process").
						AddStaticElement(metricName),
					Config:      cfg,
					Description: label.description,
					Unit:        label.unit,
				})
		}
	}

	return metricTypes, nil
}

// GetConfigPolicy returns config policy
func (procPlg *procPlugin) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{pluginVendor, fs, PluginName}, "proc_path", false, plugin.SetDefaultString("/proc"))
	return *policy, nil
}

// helper structure used to avoid duplicate metrics
type metricKey struct {
	Proc, Stat string
}

// CollectMetrics retrieves values for given metrics types
func (procPlg *procPlugin) CollectMetrics(metricTypes []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	stateCount := map[string]uint64{}

	procPath, err := metricTypes[0].Config.GetString("proc_path")
	if err != nil {
		return nil, err
	}

	// init stateCount map with keys from States
	for _, state := range States.Values() {
		stateCount[state] = 0
	}
	// get all proc stats
	stats, err := procPlg.mc.GetStats(procPath)
	if err != nil {
		return nil, err
	}
	// calculate number of processes in each state
	for _, process := range stats {
		for _, instance := range process {
			if stateName, ok := States[instance.State]; ok {
				stateCount[stateName]++
			} else {
				return nil, fmt.Errorf("Cannot find state %s is state map", stateName)
			}
		}
	}

	// calculate metrics
	aggregated := map[string]map[string]uint64{}
	processCount := map[string]uint64{}
	for _, metricType := range metricTypes {
		ns := metricType.Namespace
		if len(ns) == 7 && ns[nsCategory].Value == "process" { // process metrics
			reqProcName := ns[nsProcName].Value
			reqProcPID := ns[nsPid].Value
			metricName := ns[nsPidMetric].Value

			// return per-process metrics and build aggregated data map
			for processName, process := range stats {
				aggregated[processName] = map[string]uint64{}

				for processPid, instance := range process {
					if processName == reqProcName || reqProcName == "*" {
						if strconv.Itoa(processPid) == reqProcPID || reqProcPID == "*" {
							procMetrics, err := setProcMetrics(instance)
							if err != nil {
								return nil, fmt.Errorf("Error setting metric data: %v", err)
							}
							nuns := append([]plugin.NamespaceElement{}, ns...)
							nuns[nsProcName] = fillNsElement(&nuns[nsProcName], processName)
							nuns[nsPid] = fillNsElement(&nuns[nsPid], strconv.Itoa(processPid))

							metric := plugin.Metric{
								Namespace:   nuns,
								Data:        procMetrics[metricName],
								Timestamp:   time.Now(),
								Unit:        metricNames[metricName].unit,
								Description: metricNames[metricName].description,
							}
							metrics = append(metrics, metric)
						}

						if reqProcPID == "all" {
							procMetrics, err := setProcMetrics(instance)
							if err != nil {
								return nil, fmt.Errorf("Error setting metric data: %v", err)
							}
							for procMetricName, val := range procMetrics {
								if valInt, ok := val.(uint64); ok {
									aggregated[processName][procMetricName] += valInt
								}
							}
						}
					}
				}
			}

			// return aggregated metrics
			if reqProcPID == "all" {
				for aggrProcessName, aggrMetric := range aggregated {
					for aggrMetricName, aggrData := range aggrMetric {
						if reqProcName == aggrProcessName || reqProcName == "*" {
							if metricName == aggrMetricName {
								nuns := append([]plugin.NamespaceElement{}, ns...)
								nuns[nsProcName] = fillNsElement(&nuns[nsProcName], aggrProcessName)
								nuns[nsPid] = fillNsElement(&nuns[nsPid], "all")
								metrics = append(metrics, prepareMetric(nuns, metricName, aggrData))
							}
						}
					}
				}
			}
		} else if len(ns) == 6 && ns[nsCategory].Value == "process" { // process count
			reqProcName := ns[nsProcName].Value
			metricName := ns[nsPsCount].Value

			for processName, process := range stats {
				processCount[processName] = uint64(len(process))
			}

			// return process count metric
			for processName, processCount := range processCount {
				if reqProcName == processName || reqProcName == "*" {
					nuns := append([]plugin.NamespaceElement{}, ns...)
					nuns[nsProcName] = fillNsElement(&nuns[nsProcName], processName)
					nuns[nsPsCount] = fillNsElement(&nuns[nsPsCount], "ps_count")
					metrics = append(metrics, prepareMetric(nuns, metricName, processCount))
				}
			}
		} else if len(ns) == 5 && ns[nsCategory].Value == "state" { // globally aggregated process states
			metricName := ns[nsStateName].Value

			for stateName, val := range stateCount {
				if metricName == stateName {
					nuns := append([]plugin.NamespaceElement{}, ns...)
					metrics = append(metrics, prepareMetric(nuns, metricName, val))
				}
			}
		} else {
			return nil, fmt.Errorf("Bad namespace: %s", strings.Join(ns.Strings(), "/"))
		}
	}

	return metrics, nil
}

func setProcMetrics(instance Proc) (map[string]interface{}, error) {
	var procMetrics = make(map[string]interface{})

	if len(instance.Stat) < 29 {
		return nil, fmt.Errorf("Process instance stat data not available or broken")
	}

	vm, err := strconv.ParseUint(string(instance.Stat[22]), 10, 64)
	if err != nil {
		return nil, err
	}
	procMetrics["ps_vm"] = vm

	rss, err := strconv.ParseUint(string(instance.Stat[23]), 10, 64)
	if err != nil {
		return nil, err
	}
	procMetrics["ps_rss"] = rss

	procMetrics["ps_data"] = instance.VmData
	procMetrics["ps_code"] = instance.VmCode
	procMetrics["ps_cmdline"] = instance.CmdLine

	stack1, err := strconv.ParseUint(string(instance.Stat[27]), 10, 64)
	if err != nil {
		return nil, err
	}
	stack2, err := strconv.ParseUint(string(instance.Stat[28]), 10, 64)
	if err != nil {
		return nil, err
	}

	// to avoid overload
	if stack1 > stack2 {
		procMetrics["ps_stacksize"] = stack1 - stack2
	} else {
		procMetrics["ps_stacksize"] = stack2 - stack1
	}

	utime, err := strconv.ParseUint(string(instance.Stat[13]), 10, 64)
	if err != nil {
		return nil, err
	}
	procMetrics["ps_cputime_user"] = utime

	stime, err := strconv.ParseUint(string(instance.Stat[14]), 10, 64)
	if err != nil {
		return nil, err
	}
	procMetrics["ps_cputime_system"] = stime

	minflt, err := strconv.ParseUint(string(instance.Stat[9]), 10, 64)
	if err != nil {
		return nil, err
	}
	procMetrics["ps_pagefaults_min"] = minflt

	majflt, err := strconv.ParseUint(string(instance.Stat[11]), 10, 64)
	if err != nil {
		return nil, err
	}
	procMetrics["ps_pagefaults_maj"] = majflt

	procMetrics["ps_disk_octets_rchar"] = instance.Io["rchar"]
	procMetrics["ps_disk_octets_wchar"] = instance.Io["wchar"]
	procMetrics["ps_disk_ops_syscr"] = instance.Io["syscr"]
	procMetrics["ps_disk_ops_syscw"] = instance.Io["syscw"]

	return procMetrics, nil
}

func fillNsElement(element *plugin.NamespaceElement, value string) plugin.NamespaceElement {
	return plugin.NamespaceElement{Value: value, Description: element.Description, Name: element.Name}
}

func prepareMetric(ns []plugin.NamespaceElement, metricName string, data interface{}) plugin.Metric {
	return plugin.Metric{
		Namespace:   ns,
		Data:        data,
		Timestamp:   time.Now(),
		Unit:        metricNames[metricName].unit,
		Description: metricNames[metricName].description,
	}
}

// procPlugin holds host name and reference to metricCollector which has method of GetStats()
type procPlugin struct {
	host string
	mc   metricCollector
}

type label struct {
	description string
	unit        string
	category    string
}
