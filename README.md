# mackerel-plugin-linux-usage

mackerel metric plugin for linux usage. CPU usage (max 100%), Load average per cores, Number of processes and TCP

## Usage

```
Usage:
  mackerel-plugin-linux-usage [OPTIONS]

Application Options:
  -v, --version            Show version

Help Options:
  -h, --help               Show this help message
```

```
$ ./mackerel-plugin-linux-usage
linux-usage.cpu.guest_nice      0       1672116689
linux-usage.cpu.guest   0       1672116689
linux-usage.cpu.steal   0.004797        1672116689
linux-usage.cpu.softirq 0.155898        1672116689
linux-usage.cpu.irq     0.088742        1672116689
linux-usage.cpu.iowait  0.028781        1672116689
linux-usage.cpu.idle    98.455413       1672116689
linux-usage.cpu.system  0.765098        1672116689
linux-usage.cpu.nice    0       1672116689
linux-usage.cpu.user    0.501271        1672116689
linux-usage.loadavg.loadavg1    0       1672116689
linux-usage.loadavg.loadavg5    0       1672116689
linux-usage.loadavg.loadavg15   0       1672116689
linux-usage.process.all 86      1672116689
linux-usage.process.running     2       1672116689
linux-usage.tcp-opens.active    0       1672116689
linux-usage.tcp-opens.passive   8.325359        1672116689
linux-usage.tcp-listen.overflows        0       1672116689
linux-usage.tcp-listen.drops    0       1672116689
```

## Install

Please download release page or `mkr plugin install monitoring-forge/mackerel-plugin-linux-usage`.