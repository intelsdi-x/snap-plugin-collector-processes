// +build linux

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
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	mockPath = "./mocktest"
	mockPid  = []int{1, 12, 345, 6070, 80900}

	// mocked content of proc/<pid>/stat
	mockFileStatCont = []byte(`21926 (mockProcName) R 9018 9018 3635 34817 9018 4243788 121 0 0 0 0 0 0 0 20 0 1
								0 717086134 0 0 18446744073709551615 0 0 0 0 0 0 0 4 65536 18446744071579322839
								0 0 17 5 0 0 0 0 0 0 0 0 0 0 0 0 65280`)

	// mocked content of proc/<pid>/cmdline
	mockFileCmdlineCont = []byte(`/usr/lib/systemd/systemd-hostnamed^@`)

	// mocked content of proc/<pid>/status
	mockFileStatusCont = []byte(`
							Name:   mockProcName
							State:  R (running)
							VmData: 100
							VmExe:	100
							VmLib:	100
							Tgid:   21926
							Ngid:   0
							Pid:    21926
							PPid:   9018
							TracerPid:      0								 						 
							Uid:    0       0       0       0
							Gid:    0       0       0       0
									FDSize: 0
									Groups: 0
									Threads:        1
									SigQ:   1/127371
									SigPnd: 0000000000000000
									ShdPnd: 0000000000000000
									SigBlk: 0000000000000000
									SigIgn: 0000000000000004
									SigCgt: 0000000180010000
									CapInh: 0000000000000000
									CapPrm: 0000001fffffffff
									CapEff: 0000001fffffffff
									CapBnd: 0000001fffffffff
									Seccomp:        0
									Cpus_allowed:   ff
									Cpus_allowed_list:      0-7
									Mems_allowed:   00000000
									Mems_allowed_list:      0
									voluntary_ctxt_switches:        1
									nonvoluntary_ctxt_switches:     0
								`)

	// mocked content of proc/<pid>/io
	mockFileIoCont = []byte(`rchar: 10
							wchar: 20
							syscr: 30
							syscw: 40
							read_bytes: 50
							write_bytes: 60
							cancelled_write_bytes: 0
						`)
)

func TestGetStats(t *testing.T) {
	procPath = mockPath
	dut := &procStatsCollector{}

	Convey("new proc stats collector", t, func() {
		So(dut, ShouldNotBeNil)
	})

	Convey("when procfs directory does not exist", t, func() {
		deleteMockFiles()
		results, err := dut.GetStats()

		So(err, ShouldNotBeNil)
		So(results, ShouldBeEmpty)
	})

	Convey("when none process exist", t, func() {
		deleteMockFiles()
		os.Mkdir(mockPath, os.ModePerm)
		results, err := dut.GetStats()

		So(results, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("when some of processes files are not available", t, func() {
		files := []string{"/stat", "/cmdline", "/io", "/status"}

		for _, fileName := range files {
			createMockFiles()
			fileToRemove := mockPath + "/" + strconv.Itoa(mockPid[0]) + fileName
			os.Remove(fileToRemove)
			results, err := dut.GetStats()

			So(err, ShouldNotBeNil)
			So(results, ShouldBeEmpty)
		}
	})

	Convey("when proc files are available", t, func() {

		Convey("when proccess is not in a zombie state", func() {
			createMockFiles()
			results, err := dut.GetStats()

			So(err, ShouldBeNil)
			So(results, ShouldNotBeEmpty)

			for procName, instances := range results {

				So(procName, ShouldResemble, "mockProcName")

				for _, instance := range instances {

					// PID should be equal to one oth mocked PID
					So(mockPid, ShouldContain, instance.Pid)

					// running state in mocked content of proc/<pid>/status
					So(instance.State, ShouldEqual, "R")
					So(instance.VmData, ShouldEqual, 100*1024)
					So(instance.VmCode, ShouldEqual, (100+100)*1024) // equal to (VmExe+VMLib)*1024

					So(instance.CmdLine, ShouldResemble, string(mockFileCmdlineCont))

					So(instance.Io["cancelled_write_bytes"], ShouldEqual, 0)
					So(instance.Io["rchar"], ShouldEqual, 10)
					So(instance.Io["wchar"], ShouldEqual, 20)
					So(instance.Io["syscr"], ShouldEqual, 30)
					So(instance.Io["syscw"], ShouldEqual, 40)
					So(instance.Io["read_bytes"], ShouldEqual, 50)
					So(instance.Io["write_bytes"], ShouldEqual, 60)
				}
			}
		})

		Convey("when process is a zombie", func() {
			// change status in mockFileStatCont to zombie (Z)
			mockFileStatCont = []byte(strings.Replace(string(mockFileStatCont), " R ", " Z ", 1))

			createMockFiles()
			results, err := dut.GetStats()

			So(err, ShouldBeNil)
			So(results, ShouldNotBeEmpty)

			for _, instances := range results {

				for _, instance := range instances {

					// running state in mocked content of proc/<pid>/status
					So(instance.State, ShouldEqual, "Z")

					// no VMData and VmCode for zombie processes
					So(instance.VmData, ShouldEqual, 0)
					So(instance.VmCode, ShouldEqual, 0)
				}
			}
		})
	})

	deleteMockFiles()
}

func createMockFiles() {
	deleteMockFiles()
	os.Mkdir(mockPath, os.ModePerm)

	for _, pid := range mockPid {
		dir := mockPath + "/" + strconv.Itoa(pid)

		os.Mkdir(dir, os.ModePerm)

		f, _ := os.Create(dir + "/stat")
		f.Write(mockFileStatCont)

		f, _ = os.Create(dir + "/cmdline")
		f.Write(mockFileCmdlineCont)

		f, _ = os.Create(dir + "/status")
		f.Write(mockFileStatusCont)

		f, _ = os.Create(dir + "/io")
		f.Write(mockFileIoCont)
	}
}

func deleteMockFiles() {
	os.RemoveAll(mockPath)
}
