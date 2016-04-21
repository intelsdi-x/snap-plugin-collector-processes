# snap plugin collector - processes

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
/intel/procfs/processes/running | uint64 | Number of processes in state running
/intel/procfs/processes/sleeping | uint64 | Number of processes in state sleeping
/intel/procfs/processes/waiting | uint64 | Number of processes in state waiting
/intel/procfs/processes/zombie | uint64 | Number of processes in state zombie
/intel/procfs/processes/stopped | uint64 | Number of processes in state stopped
/intel/procfs/processes/tracing | uint64 | Number of processes in state tracing
/intel/procfs/processes/dead | uint64 | Number of processes in state dead
/intel/procfs/processes/wakekill | uint64 | Number of processes in state wakekill
/intel/procfs/processes/waking | uint64 | Number of processes in state waking
/intel/procfs/processes/parked | uint64 | Number of processes in state parked
/intel/procfs/processes/\<process_name\>/ps_vm | uint64 | Virtual memory size in bytes 
/intel/procfs/processes/\<process_name\>/ps_rss | uint64 | Resident Set Size: number of pages the process has in real memory
/intel/procfs/processes/\<process_name\>/ps_data | uint64 | Size of data segments
/intel/procfs/processes/\<process_name\>/ps_code | uint64 | Size of text segments
/intel/procfs/processes/\<process_name\>/ps_stacksize | uint64 | Stack size
/intel/procfs/processes/\<process_name\>/ps_cputime_user | uint64 | Amount of time that this process has been scheduled in user mode
/intel/procfs/processes/\<process_name\>/ps_cputime_system | uint64 | Amount of time that this process has been scheduled in kernel mode
/intel/procfs/processes/\<process_name\>/ps_pagefaults_min | uint64 | The number of minor faults the process has made
/intel/procfs/processes/\<process_name\>/ps_pagefaults_maj | uint64 | The number of major faults the process has made
/intel/procfs/processes/\<process_name\>/ps_disk_ops_syscr | uint64 | Attempt to count the number of read I/O operations
/intel/procfs/processes/\<process_name\>/ps_disk_ops_syscw | uint64 | Attempt to count the number of write I/O operations
/intel/procfs/processes/\<process_name\>/ps_disk_octets_rchar | uint64 | The number of bytes which this task has caused to be read from storage
/intel/procfs/processes/\<process_name\>/ps_disk_octets_wchar | uint64 | The number of bytes which this task has caused, or shall cause to be written to disk
/intel/procfs/processes/\<process_name\>/ps_count | uint64 | Number of process instances