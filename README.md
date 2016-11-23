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
$ snaptel metric list
NAMESPACE                                              VERSIONS
/intel/procfs/processes/running                        7
/intel/procfs/processes/sleeping                       7
/intel/procfs/processes/waiting                        7
/intel/procfs/processes/zombie                         7
/intel/procfs/processes/stopped                        7
/intel/procfs/processes/tracing                        7
/intel/procfs/processes/dead                           7
/intel/procfs/processes/wakekill                       7
/intel/procfs/processes/waking                         7
/intel/procfs/processes/parked                         7
/intel/procfs/processes/*/ps_vm                        7
/intel/procfs/processes/*/ps_rss                       7
/intel/procfs/processes/*/ps_data                      7
/intel/procfs/processes/*/ps_code                      7
/intel/procfs/processes/*/ps_stacksize                 7
/intel/procfs/processes/*/ps_cputime_user              7
/intel/procfs/processes/*/ps_cputime_system            7
/intel/procfs/processes/*/ps_pagefaults_min            7
/intel/procfs/processes/*/ps_pagefaults_maj            7
/intel/procfs/processes/*/ps_disk_ops_syscr            7
/intel/procfs/processes/*/ps_disk_ops_syscw            7
/intel/procfs/processes/*/ps_disk_octets_rchar         7
/intel/procfs/processes/*/ps_disk_octets_wchar         7
/intel/procfs/processes/*/ps_count                     7
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
