# snap plugin collector - processes

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_cmdline | string | Process command line with arguments
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_code | uint64 | Size of text segment (bytes)
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_cputime_system | uint64 | Amount of time that this process has been scheduled in kernel mode (in jiff)
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_cputime_user | uint64 | Amount of time that this process has been scheduled in user mode (in jiff)
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_data | uint64 | Size of data segments (in bytes)
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_octets_rchar | uint64 | The number of bytes which this task has caused to be read from storage (in bytes)
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_octets_wchar | uint64 | The number of bytes which this task has caused, or shall cause to be written to disk (in bytes)
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_ops_syscr | uint64 | Attempt to count the number of read I/O operations
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_disk_ops_syscw | uint64 | Attempt to count the number of write I/O operations
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_pagefaults_maj | uint64 | The number of major faults the process has made
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_pagefaults_min | uint64 | The number of minor faults the process has made
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_rss | uint64 | Resident Set Size: number of pages the process has in real memory
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_stacksize | uint64 | Stack size (in bytes)
/intel/procfs/processes/process/[process_name]/[process_pid]/ps_vm | uint64 | Virtual memory size in bytes (in bytes)
/intel/procfs/processes/process/[process_name]/all/ps_code | uint64 | Size of text segment (in bytes)
/intel/procfs/processes/process/[process_name]/all/ps_cputime_system | uint64 | Amount of time that this process has been scheduled in kernel mode (in jiff)
/intel/procfs/processes/process/[process_name]/all/ps_cputime_user | uint64 | Amount of time that this process has been scheduled in user mode (in jiff)
/intel/procfs/processes/process/[process_name]/all/ps_data | uint64 | Size of data segments (in bytes)
/intel/procfs/processes/process/[process_name]/all/ps_disk_octets_rchar | uint64 | The number of bytes which this task has caused to be read from storage (in bytes)
/intel/procfs/processes/process/[process_name]/all/ps_disk_octets_wchar | uint64 | The number of bytes which this task has caused, or shall cause to be written to disk (in bytes)
/intel/procfs/processes/process/[process_name]/all/ps_disk_ops_syscr | uint64 | Attempt to count the number of read I/O operations
/intel/procfs/processes/process/[process_name]/all/ps_disk_ops_syscw | uint64 | Attempt to count the number of write I/O operations
/intel/procfs/processes/process/[process_name]/all/ps_pagefaults_maj | uint64 | The number of major faults the process has made
/intel/procfs/processes/process/[process_name]/all/ps_pagefaults_min | uint64 | The number of minor faults the process has made
/intel/procfs/processes/process/[process_name]/all/ps_rss | uint64 | Resident Set Size: number of pages the process has in real memory
/intel/procfs/processes/process/[process_name]/all/ps_stacksize | uint64 | Stack size (in bytes)
/intel/procfs/processes/process/[process_name]/all/ps_vm | uint64 | Virtual memory size (in bytes)
/intel/procfs/processes/process/[process_name]/ps_count | uint64 | Number of process instances
/intel/procfs/processes/state/dead | uint64 | Number of processes with 'dead' status
/intel/procfs/processes/state/parked | uint64 | Number of processes with 'parked' status
/intel/procfs/processes/state/running | uint64 | Number of processes with 'running' status
/intel/procfs/processes/state/sleeping | uint64 | Number of processes with 'sleeping' status
/intel/procfs/processes/state/stopped | uint64 | Number of processes with 'stopped' status
/intel/procfs/processes/state/tracing | uint64 | Number of processes with 'tracing' status
/intel/procfs/processes/state/waiting | uint64 | Number of processes with 'waiting' status
/intel/procfs/processes/state/wakekill | uint64 | Number of processes with 'wakekill' status
/intel/procfs/processes/state/waking | uint64 | Number of processes with 'waking' status
/intel/procfs/processes/state/zombie | uint64 | Number of processes with 'zombie' status
