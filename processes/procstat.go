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
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/intelsdi-x/snap-plugin-utilities/str"
	log "github.com/sirupsen/logrus"
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
	Stat    []string
	Io      map[string]uint64
	VmData  uint64
	VmCode  uint64
}

// GetStats returns processes statistics
func (psc *procStatsCollector) GetStats(procPath string) (map[string]map[int]Proc, error) {
	// Procfs structure used in GetStats
	// /proc
	// |_ /[pid] (for example 922)
	//    |_ cmdline (process command like, for example /usr/local/bin/snapteld -t 0 -l 1)
	//    |_ io (I/O information about the process)
	//    |_ stat (Status information about the process)
	//    |_ status (Provides much of the information in /proc/[pid]/stat and
	//               /proc/[pid]/statm in a format that's easier for humans to
	//               parse)
	//
	// For more details, check here:
	// http://man7.org/linux/man-pages/man5/proc.5.html

	files, err := ioutil.ReadDir(procPath)
	if err != nil {
		return nil, err
	}
	procs := map[string]map[int]Proc{}
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

			procStatFields := strings.Fields(string(procStat))
			if len(procStatFields) < 3 {
				return nil, fmt.Errorf("Cannot retrieve process state")
			}
			procState := procStatFields[2]

			// TODO: gather task status data /proc/<pid>/task
			pc := Proc{
				Pid:     pid,
				State:   procState,
				Stat:    procStatFields,
				CmdLine: strings.Replace(string(procCmdLine), "\x00", " ", -1),
				Io:      procIo,
				VmData:  vmData,
				VmCode:  vmCode,
			}
			// procName is process name extracted from command line path
			procPath := strings.Split(pc.CmdLine, " ")[0]
			procName := ""
			if procPath != "" {
				procName = filepath.Base(procPath)
			} else {
				// Kernel processes - no command line
				procName = removeUnwantedChars(procStatFields[1])
			}
			if procName == "" {
				return nil, fmt.Errorf("Cannot retrieve process name")
			}

			if procs[procName] == nil {
				procs[procName] = map[int]Proc{}
			}
			procs[procName][pid] = pc
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

type metricCollector interface {
	GetStats(procPath string) (map[string]map[int]Proc, error)
}

type unwanted struct {
	char string
	repl string
}

// for mocking
type procStatsCollector struct{}
