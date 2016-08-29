# snap plugin collector - processes

snap plugin for collects information about process states grouped by name. Additionally it provides metrics for each running process, also grouped by name. 

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

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-processes/blob/master/README.md#examples).

## Documentation

### Collected Metrics
List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-processes/blob/master/METRICS.md).

### Examples

Example running snap-plugin-collector-processes plugin and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```

Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```
In another terminal window:

Load snap-plugin-collector-processes plugin:
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-processes
```
Load mock-file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-plugin-publisher-mock-file
```
See available metrics for your system:
```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file to use snap-plugin-collector-processes plugin (exemplary files in [examples/tasks/] (examples/tasks/)):
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
                "/intel/procfs/processes/*/ps_disk_ops_syscr": {},
                "/intel/procfs/processes/*/ps_disk_ops_syscw": {},
                "/intel/procfs/processes/running": {},
                "/intel/procfs/processes/stopped": {},
                "/intel/procfs/processes/waiting": {}
            },
            "publish": [
                {
                    "plugin_name": "mock-file",
                    "config": {
                        "file": "/tmp/published_processes.log"
                    }
                }
            ],
            "config": null
        }
    }
}

```

Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/processes-file.json
```

If you would like to collect all metrics exposed by this plugin, set `/intel/procfs/processes/*` as a metric to collect in task manifest.

### Roadmap
As we launch this plugin, we have a few items in mind for the next release:

- [  ] Gathering task status from `"/proc/<pid>/task"`

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-processes/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-processes/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Marcin Krolik](https://github.com/marcin-krolik)
* Co-author: [Izabella Raulin](https://github.com/IzabellaRaulin)
