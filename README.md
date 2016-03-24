# snap-plugin-collector-processes

snap plugin for collects information about process states grouped by name. Additionally it provides metrics for each running process, also grouped by name. 

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

 Plugin collects specified metrics in-band on OS level

### System Requirements

 - Linux system with proc filesystem

### Installation
#### Download processes plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-processes
Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-processes
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description (optional)
----------|-----------|-----------------------
/intel/procfs/processes/running | | Number of processes in state running
/intel/procfs/processes/sleeping | | Number of processes in state sleeping
/intel/procfs/processes/waiting | | Number of processes in state waiting
/intel/procfs/processes/zombie | | Number of processes in state zombie
/intel/procfs/processes/stopped | | Number of processes in state stopped
/intel/procfs/processes/tracing | | Number of processes in state tracing
/intel/procfs/processes/dead | | Number of processes in state dead
/intel/procfs/processes/wakekill | | Number of processes in state wakekill
/intel/procfs/processes/waking | | Number of processes in state waking
/intel/procfs/processes/parked | | Number of processes in state parked
/intel/procfs/processes/\<proces_name\>/ps_vm | | Virtual memory size in bytes 
/intel/procfs/processes/\<proces_name\>/ps_rss | | Resident Set Size: number of pages the process has in real memory
/intel/procfs/processes/\<proces_name\>/ps_data | | Size of data segments
/intel/procfs/processes/\<proces_name\>/ps_code | | Size of text segments
/intel/procfs/processes/\<proces_name\>/ps_stacksize | | Stack size
/intel/procfs/processes/\<proces_name\>/ps_cputime_user | | Amount of time that this process has been scheduled in user mode
/intel/procfs/processes/\<proces_name\>/ps_cputime_system | | Amount of time that this process has been scheduled in kernel mode
/intel/procfs/processes/\<proces_name\>/ps_pagefaults_min | | The number of minor faults the process has made
/intel/procfs/processes/\<proces_name\>/ps_pagefaults_maj | | The number of major faults the process has made
/intel/procfs/processes/\<proces_name\>/ps_disk_ops_syscr | | Attempt to count the number of read I/O operations
/intel/procfs/processes/\<proces_name\>/ps_disk_ops_syscw | | Attempt to count the number of write I/O operations
/intel/procfs/processes/\<proces_name\>/ps_disk_octets_rchar | | The number of bytes which this task has caused to be read from storage
/intel/procfs/processes/\<proces_name\>/ps_disk_octets_wchar | | The number of bytes which this task has caused, or shall cause to be written to disk
/intel/procfs/processes/\<proces_name\>/ps_count | | Number of process instances

### Examples
Example task manifest to use processes plugin:
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/procfs/processes/*": {},
                "/intel/procfs/processes/running": {},
                "/intel/procfs/processes/sleeping": {},
                "/intel/procfs/processes/zombie": {}
            },
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_processes"
                    }
                }
            ],
            "config": null
        }
    }
}

```


### Roadmap

- gather task status `"/proc/<pid>/task"`

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-kairosdb/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Marcin Krolik](https://github.com/marcin-krolik)
* Co-author: [Izabella Raulin](https://github.com/IzabellaRaulin)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.