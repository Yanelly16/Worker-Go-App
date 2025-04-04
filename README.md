# Worker-Go-App
## Overview
This is a high-performance TCP port scanner built with Go. It features concurrent scanning, banner grabbing, and supports both human-readable and JSON output formats. The scanner efficiently checks port statuses across multiple targets with configurable workers and timeout settings.

##  Features ✨

✅ Concurrent scanning with adjustable workers  
✅ Multiple target support (comma-separated list)  
✅ Banner grabbing for service identification  
✅ Configurable port ranges and timeouts  
✅ Dual output formats (human-readable and JSON)  
✅ Progress indicators during long scans  
✅ Specific port scanning capability  

## **Setup Instructions** 📦

### **Prerequisites**
- Go 1.16+ installed
- Network access to target systems

### **Steps**
- cd into the directory 
  cd worker-go-app

### **Usage Guide**
- Basic the application
  go build -o port-scanner
- Basic Scanning
./port-scanner -targets= scanme.nmap.org

## **Sample Output**
Port 22 is open - SSH-2.0-OpenSSH_8.2
Port 80 is open - HTTP/1.1 200 OK

Scan Summary:
Target: example.com
Ports scanned: 1024
Open ports: 2
Time taken: 5.26s

## **JSON Output**

{
  "target": "example.com",
  "ports_scanned": 1024,
  "open_ports": 2,
  "time_taken": "5.26s",
  "results": [
    {
      "port": 22,
      "state": "open",
      "banner": "SSH-2.0-OpenSSH_8.2"
    },
    {
      "port": 80,
      "state": "open",
      "banner": "HTTP/1.1 200 OK"
    }
  ]
}

## **Demo Video**
https://youtu.be/2M34Jkyoj-A

