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
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-utilities/str"
)

const (
	procStat   = "stat"
	procStatus = "status"
	procCmd    = "cmdline"
	procIO     = "io"
)

var (
	// States contains possible states of processes
	States = str.StringMap{
		"R": "running",
		"S": "sleeping",
		"D": "waiting",
		"Z": "zombie",
		"T": "stopped",
		"t": "tracing",
		"X": "dead",
		"K": "wakekill",
		"W": "waking",
		"P": "parked",
	}
)

// Proc holds processes statistics
type Proc struct {
	Pid     int
	State   string
	CmdLine string
	Cmd     string
	Stat    []string
	Io      map[string]uint64
	VmData  uint64
	VmCode  uint64
}

// GetStats returns processes statistics
func (psc *procStatsCollector) GetStats(procPath string, isp bool) ([]Proc, error) {
	files, err := ioutil.ReadDir(procPath)
	//log.SetOutput(os.Stderr)
	if err != nil {
		return nil, err
	}
	procs := make([]Proc, 0)
	for _, file := range files {

		// process only PID sub dirs
		if pid, err := strconv.Atoi(file.Name()); err == nil {
			// get proc/<pid>/stat data
			fstat := filepath.Join(procPath, file.Name(), procStat)
			procStat, err := ioutil.ReadFile(fstat)
			if err != nil {
				log.WithFields(log.Fields{
					"pid":   pid,
					"file":  fstat,
					"error": err,
				}).Errorf("Cannot get status information about the process")
				continue
			}
			// get proc/<pid>/cmdline data
			fcmd := filepath.Join(procPath, file.Name(), procCmd)
			procCmdLine, err := ioutil.ReadFile(fcmd)
			if err != nil {
				log.WithFields(log.Fields{
					"pid":   pid,
					"file":  fcmd,
					"error": err,
				}).Errorf("Cannot get command line for the process")
				continue
			}
			// get proc/<pid>/io data
			fio := filepath.Join(procPath, file.Name(), procIO)
			procIo, err := read2Map(fio)
			if err != nil {
				log.WithFields(log.Fields{
					"pid":   pid,
					"file":  fio,
					"error": err,
				}).Errorf("Cannot get I/O statistics for the process")
				continue
			}
			// get proc/<pid>/status data
			var pStatus map[string]uint64
			var vmData, vmCode uint64
			// special case for zombie
			state := strings.Fields(string(procStat))[2]
			if state == "Z" {
				vmData = 0
				vmCode = 0
			} else {
				fstatus := filepath.Join(procPath, file.Name(), procStatus)
				pStatus, err = read2Map(fstatus)
				if err != nil {
					log.WithFields(log.Fields{
						"pid":   pid,
						"file":  fstatus,
						"error": err,
					}).Errorf("Cannot get status information for the process")
					continue
				}
				vmData = pStatus["VmData"] * 1024
				vmCode = (pStatus["VmExe"] + pStatus["VmLib"]) * 1024
			}
			// TODO: gather task status data /proc/<pid>/task
			cmdLine := string(procCmdLine)
			//			cmdPath := strings.Split(strings.Split(cmdLine, "\x00")[0], "/")
			cmdPath := strings.Split(cmdLine, "\x00")[0]
			var cmd string
			if strings.Contains(cmdPath, "[") && strings.Contains(cmdPath, "]") {
				cmd = strings.Split(cmdPath, " ")[0]
			} else {
				s := strings.Split(cmdPath, "/")
				cmd = s[len(s)-1]
			}
			if isp || len(cmd) > 0 {
				pc := Proc{
					Pid:     pid,
					State:   strings.Fields(string(procStat))[2],
					Stat:    strings.Fields(string(procStat)),
					CmdLine: cmdLine,
					Cmd:     cmd,
					Io:      procIo,
					VmData:  vmData,
					VmCode:  vmCode,
				}
				if len(pc.Cmd) == 0 {
					pc.Cmd = strings.Fields(string(procStat))[1]
				}
				procs = append(procs, pc)
			}
		}
	}
	return procs, nil
}

// readToMap retrieves statistics from file specified by filename and returns its (name, value) as a map
func read2Map(fileName string) (map[string]uint64, error) {
	stats := map[string]uint64{}
	status, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(status), "\n") {
		if line == "" {
			continue
		}

		data := strings.Fields(line)
		if len(data) < 2 {
			continue
		}

		name := data[0]
		last := len(name) - 1
		if string(name[last]) == ":" {
			name = name[:last]
		}

		value, err := strconv.ParseUint(data[1], 10, 64)

		if err != nil {
			continue
		}

		stats[name] = value
	}
	return stats, nil
}

type metricCollector interface {
	GetStats(procPath string, isp bool) ([]Proc, error)
}

type unwanted struct {
	char string
	repl string
}

// for mocking
type procStatsCollector struct{}
