# Snap plugin collector - processes

This plugin for [Snap Telemetry Framework](http://github.com/intelsdi-x/snap) collects information about process states grouped by name. Additionally it provides metrics for each running process, also grouped by name.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements

 - Linux (with proc filesystem)

### Installation

#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-processes/releasess) page. Download the plugin from the latest release and load it into `snapteld` (`/opt/snap/plugins` is the default location for Snap packages).
#### To build the plugin binary:
Use https://github.com/intelsdi-x/snap-plugin-collector-processes or your fork as repo.

Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-processes
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-processes/blob/master/README.md#examples).

Configuration parameters:

- `proc_path`: path to procfs (default: `/proc`)

## Documentation

This collector gathers metrics from proc file system. The configuration `proc_path` determines where the plugin obtains these metrics, with a default setting of `/proc`. This setting is only required to obtain data from a docker container that mounts the host `/proc` in an alternative path.

### Collected Metrics
List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-processes/blob/master/METRICS.md).

### Examples

Example of running Snap processes collector and writing data to file.

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `snapteld -l 1 -t 0 &`

Download and load snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-processes/latest/linux/x86_64/snap-plugin-collector-processes
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-processes
$ snaptel plugin load snap-plugin-publisher-file
```
See available metrics for your system:
```
$ snaptel metric list --verbose 
NAMESPACE                                                                            VERSION    UNIT    DESCRIPTION
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_cmdline              8                  Process command line with arguments
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_code                 8          B       Size of text segment
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_cputime_system       8          Jiff    Amount of time that this process has been scheduled in kernel mode
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_cputime_user         8          Jiff    Amount of time that this process has been scheduled in user mode
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_data                 8          B       Size of data segments
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_octets_rchar    8          B       The number of bytes which this task has caused to be read from storage
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_octets_wchar    8          B       The number of bytes which this task has caused, or shall cause to be written to disk
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_ops_syscr       8                  Attempt to count the number of read I/O operations
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_ops_syscw       8                  Attempt to count the number of write I/O operations
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_pagefaults_maj       8                  The number of major faults the process has made
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_pagefaults_min       8                  The number of minor faults the process has made
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_rss                  8                  Resident Set Size: number of pages the process has in real memory
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_stacksize            8          B       Stack size
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_vm                   8          B       Virtual memory size in bytes
/intel/procfs/processes/process/[process_name]/all/ps_code                           8          B       Size of text segment
/intel/procfs/processes/process/[process_name]/all/ps_cputime_system                 8          Jiff    Amount of time that this process has been scheduled in kernel mode
/intel/procfs/processes/process/[process_name]/all/ps_cputime_user                   8          Jiff    Amount of time that this process has been scheduled in user mode
/intel/procfs/processes/process/[process_name]/all/ps_data                           8          B       Size of data segments
/intel/procfs/processes/process/[process_name]/all/ps_disk_octets_rchar              8          B       The number of bytes which this task has caused to be read from storage
/intel/procfs/processes/process/[process_name]/all/ps_disk_octets_wchar              8          B       The number of bytes which this task has caused, or shall cause to be written to disk
/intel/procfs/processes/process/[process_name]/all/ps_disk_ops_syscr                 8                  Attempt to count the number of read I/O operations
/intel/procfs/processes/process/[process_name]/all/ps_disk_ops_syscw                 8                  Attempt to count the number of write I/O operations
/intel/procfs/processes/process/[process_name]/all/ps_pagefaults_maj                 8                  The number of major faults the process has made
/intel/procfs/processes/process/[process_name]/all/ps_pagefaults_min                 8                  The number of minor faults the process has made
/intel/procfs/processes/process/[process_name]/all/ps_rss                            8                  Resident Set Size: number of pages the process has in real memory
/intel/procfs/processes/process/[process_name]/all/ps_stacksize                      8          B       Stack size
/intel/procfs/processes/process/[process_name]/all/ps_vm                             8          B       Virtual memory size in bytes
/intel/procfs/processes/process/[process_name]/ps_count                              8                  Number of process instances
/intel/procfs/processes/state/dead                                                   8                  Number of processes with 'dead' status
/intel/procfs/processes/state/parked                                                 8                  Number of processes with 'parked' status
/intel/procfs/processes/state/running                                                8                  Number of processes with 'running' status
/intel/procfs/processes/state/sleeping                                               8                  Number of processes with 'sleeping' status
/intel/procfs/processes/state/stopped                                                8                  Number of processes with 'stopped' status
/intel/procfs/processes/state/tracing                                                8                  Number of processes with 'tracing' status
/intel/procfs/processes/state/waiting                                                8                  Number of processes with 'waiting' status
/intel/procfs/processes/state/wakekill                                               8                  Number of processes with 'wakekill' status
/intel/procfs/processes/state/waking                                                 8                  Number of processes with 'waking' status
/intel/procfs/processes/state/zombie                                                 8                  Number of processes with 'zombie' status

```


Download an [example task file](https://github.com/intelsdi-x/snap-plugin-collector-processes/blob/master/examples/tasks/) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-processes/master/examples/tasks/processes-file.json
$ snaptel task create -t processes-file.json
```

If you would like to collect all metrics exposed by this plugin, set `/intel/procfs/processes/*` as a metric to collect in task manifest.

### Roadmap
As we launch this plugin, we have a few items in mind for the next release:

- [  ] Gathering task status from `"/proc/<pid>/task"`

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-processes/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-processes/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [Slack](http://slack.snap-telemetry.io).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).
## Acknowledgements

* Author: [Marcin Krolik](https://github.com/marcin-krolik)
* Co-author: [Izabella Raulin](https://github.com/IzabellaRaulin)

**Thank you!** Your contribution is incredibly important to us.
