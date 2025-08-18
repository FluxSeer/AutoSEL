
# AutoSEL

## Overview

AutoSEL automates the generation, update, and application of SELinux policies by monitoring container events in combination with container configuration details.  
It removes the need to manually load configuration files, significantly reduces repetitive operations, and improves responsiveness to configuration changes. Run it once to complete policy generation and updates. Privileged containers and resource mounts are controlled more strictly.

---

## Installation and Usage

### Prerequisites
- Go 1.18 or newer
- Docker environment
- SELinux enabled (recommended on Linux distributions with SELinux support)

### Installation
```bash
git clone https://github.com/ChongZhe001025/AutoSEL.git
cd AutoSEL
go mod tidy
go build -o autosel main.go
```

### Run
```bash
sudo ./autosel
```

---

## Project Structure

```
AutoSEL/
│
├── main.go                # Entry point
├── go.mod / go.sum        # Go modules
├── README.md
│
├── applicator/
│   └── ApplyPolicy.go     # Policy application logic
│
├── monitor/
│   └── ContainerMonitor.go # Container event monitoring
│
├── parser/
│   ├── Parse.go
│   └── auxiliary/
│       └── Auxiliary.go
│   └── config/
│       └── Parse.go
│
├── policy/
│   ├── Policy.go
│   └── capability/
│       └── PolicyForCaps.go
│   └── device/
│       └── PolicyForDevices.go
│   └── mount/
│       └── PolicyForMounts.go
│
└── SysFiles/
		├── ExportedTarFiles/
		│   └── example.tar
		└── SELinuxPolicies/
				└── example.cil
```

