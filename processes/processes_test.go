// +build linux,small

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015-2016 Intel Corporation

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
	"strconv"
	"strings"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

const (
	mockProcName = "NetworkManager"
	mockProcPid  = 815

	mockProcName2 = "fake"
	mockProcPid2  = 315

	mockProcName3 = "fake"
	mockProcPid3  = 316
)

var (
	mockProc  = makeMockProc(mockProcName, mockProcPid)
	mockProc2 = makeMockProc(mockProcName2, mockProcPid2)
	mockProc3 = makeMockProc(mockProcName3, mockProcPid3)
)

type mcMock struct {
	mock.Mock
}

func (mc *mcMock) GetStats(procPath string) (map[string]map[int]Proc, error) {
	args := mc.Called()
	var r0 map[string]map[int]Proc
	if args.Get(0) != nil {
		r0 = args.Get(0).(map[string]map[int]Proc)
	}
	return r0, args.Error(1)
}

func TestGetConfigPolicy(t *testing.T) {

	Convey("normal case", t, func() {
		procPlugin := New()

		Convey("new processess collector", func() {
			So(procPlugin, ShouldNotBeNil)
		})

		So(func() { procPlugin.GetConfigPolicy() }, ShouldNotPanic)
		_, err := procPlugin.GetConfigPolicy()
		So(err, ShouldBeNil)
	})
}

func TestGetMetricTypes(t *testing.T) {

	var cfg plugin.Config

	Convey("get metric types successfully", t, func() {
		procPlugin := New()
		Convey("new processes collector", func() {
			So(procPlugin, ShouldNotBeNil)
		})

		So(func() { procPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
		results, err := procPlugin.GetMetricTypes(cfg)

		So(err, ShouldBeNil)
		So(results, ShouldNotBeEmpty)

		// plugin returns total of 38 metrics available, see the README.md
		So(len(results), ShouldEqual, 38)

		for _, res := range results {
			So(res.Description, ShouldNotBeBlank)
		}
	})
}

func TestCollectMetrics(t *testing.T) {

	Convey("collect metric", t, func() {
		procPlugin := New()

		cfg := plugin.Config{
			"proc_path": "/proc",
		}

		mockMts := []plugin.Metric{
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "dead"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "parked"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "running"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "sleeping"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "stopped"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "tracing"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "waiting"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "wakekill"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "waking"),
				Config:    cfg,
			},
			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "state", "zombie"),
				Config:    cfg,
			},

			plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process").
					AddDynamicElement("process_name", "name of the process").
					AddDynamicElement("process_pid", "identifier of the process").
					AddStaticElement("ps_data"),
			},
		}

		mockMts[len(mockMts)-1].Namespace[4].Value = "fake"

		Convey("new processes collector", func() {
			So(procPlugin, ShouldNotBeNil)
		})

		Convey("when attempt to get processes statistics fails with error", func() {
			mc := &mcMock{}
			procPlugin.mc = mc

			mc.On("GetStats", mock.Anything).Return(nil, errors.New("x"))

			results, err := procPlugin.CollectMetrics(mockMts)

			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		})

		Convey("when getStats() returns list of valid processes statistics", func() {
			mc := &mcMock{}
			procPlugin.mc = mc

			mc.On("GetStats").Return(map[string]map[int]Proc{
				"NetworkManager": map[int]Proc{
					mockProcPid: mockProc,
				},
				"fake": map[int]Proc{
					mockProcPid2: mockProc2,
					mockProcPid3: mockProc3,
				},
			}, nil)

			Convey("when names of collect metrics are valid", func() {
				results, err := procPlugin.CollectMetrics(mockMts)

				So(err, ShouldBeNil)
				// 12 metrics exposed - 11 metrics defined in mockMts, where last metric returns data for 2 PIDs
				So(len(results), ShouldEqual, 12)
			})

			Convey("when name of collect metric process name and pid is equal to asterisk (exposed dynamic metrics)", func() {

				results, err := procPlugin.CollectMetrics([]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process").
							AddDynamicElement("process_name", "name of the process").
							AddDynamicElement("process_pid", "identifier of the process").
							AddStaticElement("ps_data"),
						Config: cfg,
					},
				})

				So(err, ShouldBeNil)
				// 3 metrics exposed by processes
				So(len(results), ShouldEqual, 3)

				for _, r := range results {
					ns := r.Namespace
					So(mockProc.validateValue(ns[len(ns)-1].Value, r.Data.(uint64)), ShouldBeTrue)
				}
			})

			Convey("when name of collect metric process name is equal to asterisk and pid specified", func() {

				results, err := procPlugin.CollectMetrics([]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process").
							AddDynamicElement("process_name", "name of the process").
							AddStaticElement(strconv.Itoa(mockProcPid)).
							AddStaticElement("ps_data"),
						Config: cfg,
					},
				})

				So(err, ShouldBeNil)
				// 1 metrics exposed by specific process pid
				So(len(results), ShouldEqual, 1)

				for _, r := range results {
					ns := r.Namespace
					So(mockProc.validateValue(ns[len(ns)-1].Value, r.Data.(uint64)), ShouldBeTrue)
				}
			})

			Convey("when name of collect metric process name is specified and pid is equal to asterisk", func() {

				results, err := procPlugin.CollectMetrics([]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process").
							AddStaticElement("fake").
							AddDynamicElement("process_pid", "identifier of the process").
							AddStaticElement("ps_data"),
						Config: cfg,
					},
				})

				So(err, ShouldBeNil)
				// 2 metrics exposed by fake process with 2 pids
				So(len(results), ShouldEqual, 2)

				for _, r := range results {
					ns := r.Namespace
					So(mockProc.validateValue(ns[len(ns)-1].Value, r.Data.(uint64)), ShouldBeTrue)
				}
			})

			Convey("check process instance count with process name specified", func() {
				results, err := procPlugin.CollectMetrics([]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process", "fake", "ps_count"),
						Config:    cfg,
					},
					plugin.Metric{
						Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process", "NetworkManager", "ps_count"),
						Config:    cfg,
					},
				})

				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 2)

				matchedNs := 0
				for _, r := range results {
					ns := strings.Join(r.Namespace.Strings(), "/")

					switch ns {
					case "intel/procfs/processes/process/fake/ps_count":
						matchedNs++
						So(r.Data, ShouldEqual, 2)
					case "intel/procfs/processes/process/NetworkManager/ps_count":
						matchedNs++
						So(r.Data, ShouldEqual, 1)
					}
				}

				So(matchedNs, ShouldEqual, 2)
			})

			Convey("check process instance count with process name equal to asterisk", func() {
				results, err := procPlugin.CollectMetrics([]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process").
							AddDynamicElement("process_name", "name of the process").
							AddStaticElement("ps_count"),
						Config: cfg,
					},
				})

				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 2)

				matchedNs := 0
				for _, r := range results {
					ns := strings.Join(r.Namespace.Strings(), "/")

					switch ns {
					case "intel/procfs/processes/process/fake/ps_count":
						matchedNs++
						So(r.Data, ShouldEqual, 2)
					case "intel/procfs/processes/process/NetworkManager/ps_count":
						matchedNs++
						So(r.Data, ShouldEqual, 1)
					}
				}

				So(matchedNs, ShouldEqual, 2)
			})

			Convey("when names of collect metrics include asterisk", func() {
				mockMtsWithAsterisk := append(mockMts, plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process").
						AddDynamicElement("process_name", "name of the process").
						AddDynamicElement("process_pid", "identifier of the process").
						AddStaticElement("ps_data"),
				})

				results, err := procPlugin.CollectMetrics(mockMtsWithAsterisk)

				So(err, ShouldBeNil)
				// 5 dynamic metrics exposed by processes + 10 status metrics defined in mockMts
				So(len(results), ShouldEqual, 15)
			})

			Convey("when name of collect metric is invalid", func() {
				results, err := procPlugin.CollectMetrics([]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("intel", "procfs", "zombie"),
						Config:    cfg,
					},
				})

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "Bad namespace")
				So(results, ShouldBeEmpty)
			})

		})

		Convey("when getStats() returns statistics for multiple processes", func() {
			mc := &mcMock{}
			procPlugin.mc = mc
			makeSleepyName := func(num int) string {
				return fmt.Sprintf("sleepy%d", num)
			}
			statsRes := map[string]map[int]Proc{}
			for i := 0; i < 10; i++ {
				procName := makeSleepyName(i)
				statsRes[procName] = map[int]Proc{}
				statsRes[procName][i] = makeMockProc(procName, 24000+i)
			}
			numProcs := len(statsRes)
			mc.On("GetStats").Return(statsRes, nil)
			mockMtsWithAsterisk := []plugin.Metric{plugin.Metric{
				Namespace: plugin.NewNamespace("intel", "procfs", "processes", "process").
					AddDynamicElement("process_name", "name of the process").
					AddDynamicElement("process_pid", "identifier of the process").
					AddStaticElement("ps_count"),
				Config: cfg,
			}}
			numMts := len(mockMtsWithAsterisk)
			Convey("and metrics for all processes are requested in call to CollectMetrics", func() {
				results, err := procPlugin.CollectMetrics(mockMtsWithAsterisk)
				Convey("then no error should be reported", func() {
					So(err, ShouldBeNil)
				})
				Convey("so all metrics and processes should be represented in results", func() {
					So(len(results), ShouldEqual, numProcs*numMts)
				})
				Convey("so all processes should be represented in results", func() {
					expProcNames := map[string]bool{}
					actProcNames := map[string]bool{}
					for i := 0; i < numProcs; i++ {
						expProcNames[makeSleepyName(i)] = true
					}
					for _, mt := range results {
						actProcNames[mt.Namespace[4].Value] = true
					}
					So(actProcNames, ShouldResemble, expProcNames)
				})
				Convey("so all metrics should be represented in results", func() {
					expMtNames := map[string]bool{}
					actMtNames := map[string]bool{}
					for _, mt := range mockMtsWithAsterisk {
						expMtNames[mt.Namespace[6].Value] = true
					}
					for _, mt := range results {
						actMtNames[mt.Namespace[6].Value] = true
					}
					So(actMtNames, ShouldResemble, expMtNames)
				})
			})
		})

	})

}

func (mp Proc) validateValue(param string, value uint64) bool {
	ok := false

	var refValue uint64

	switch param {
	case "ps_data":
		refValue = mp.VmData
	case "ps_stacksize":
		s1, _ := strconv.ParseUint(string(mp.Stat[27]), 10, 64)
		s2, _ := strconv.ParseUint(string(mp.Stat[28]), 10, 64)

		if s1 > s2 {
			refValue = s1 - s2
		} else {
			refValue = s2 - s1
		}

	case "ps_disk_ops_syscr":
		refValue = mp.Io["syscr"]

	case "ps_disk_ops_syscw":
		refValue = mp.Io["syscw"]
	case "ps_rss":
		rssSingle, _ := strconv.ParseUint(string(mp.Stat[23]), 10, 6)
		refValue = rssSingle
	case "ps_code":
		refValue = mp.VmCode
	case "ps_cputime_system":
		stime, _ := strconv.ParseUint(string(mp.Stat[14]), 10, 64)
		refValue = stime * 10000
	case "ps_vm":
		vm, _ := strconv.ParseUint(string(mp.Stat[22]), 10, 64)
		refValue = vm
	case "ps_pagefaults_min":
		minflt, _ := strconv.ParseUint(string(mp.Stat[9]), 10, 64)
		refValue = minflt * 10000
	case "ps_pagefaults_maj":
		majflt, _ := strconv.ParseUint(string(mp.Stat[11]), 10, 64)
		refValue = majflt * 10000
	case "ps_disk_octets_rchar":
		refValue = mp.Io["rchar"]
	case "ps_disk_octets_wchar":
		refValue = mp.Io["wchar"]
	case "ps_cputime_user":
		utime, _ := strconv.ParseUint(string(mp.Stat[13]), 10, 64)
		refValue = utime * 10000
	default:
		fmt.Println("invalid metric name", param)
		return false
	} // end of switch case

	if value == refValue {
		ok = true
	} else {
		fmt.Println("invalid value for", param, "( expecting", refValue, ", is", value, ")")
	}

	return ok
}

func makeMockProc(procName string, procPid int) Proc {
	pidStr := strconv.Itoa(procPid)
	res := Proc{
		Pid:     procPid,
		State:   "S",
		CmdLine: "/usr/sbin/" + procName + " --no-daemon",

		Stat: []string{
			pidStr, "(" + procName + ")", "S", "1", pidStr, pidStr, "0", "-1", "1077960960", "3601", "513", "0", "0",
			"115", "28", "0", "0", "20", "0", "4", "331", "459870208", "2145", "18446744073709551615", "140096990736384",
			"140096992449927", "140729036690976", "140729036689856", "140096924699517", "0", "20483", "4096", "65536",
			"18446744073709551615", "0", "0", "17", "7", "0", "0", "3", "0", "0", "140096994547816", "140096994587072",
			"140097024917504", "140729036697458", "140729036697495", "140729036697495", "140729036697567", "0",
		},
		Io: map[string]uint64{
			"syscr":                 1100676,
			"syscw":                 124253,
			"read_bytes":            102400,
			"write_bytes":           0,
			"cancelled_write_bytes": 0,
			"rchar":                 260972212,
			"wchar":                 995958,
		},
		VmData: 227209216,
		VmCode: 27209216,
	}
	return res
}
