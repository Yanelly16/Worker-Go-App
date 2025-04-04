// Filename: main.go
package main

// ScanResult represents the result of scanning a single port
type ScanResult struct {
    Port   int    `json:"port"`
    State  string `json:"state"`
    Banner string `json:"banner,omitempty"`
}

// ScanSummary aggregates all scan results
type ScanSummary struct {
    Target       string       `json:"target"`
    StartPort    int          `json:"start_port,omitempty"`
    EndPort      int          `json:"end_port,omitempty"`
    PortsScanned int          `json:"ports_scanned"`
    OpenPorts    int          `json:"open_ports"`
    TimeTaken    string       `json:"time_taken"`
    Results      []ScanResult `json:"results,omitempty"`
}
var (
    targets    string
    startPort  int
    endPort    int
    workers    int
    timeout    int
    jsonOutput bool
    portsList  string
)

func init() {
    flag.StringVar(&targets, "targets", "scanme.nmap.org", "Comma-separated list of targets")
    flag.IntVar(&startPort, "start-port", 1, "Start port number")
    flag.IntVar(&endPort, "end-port", 1024, "End port number")
    flag.IntVar(&workers, "workers", 100, "Number of concurrent workers")
    flag.IntVar(&timeout, "timeout", 5, "Connection timeout in seconds")
    flag.BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
    flag.StringVar(&portsList, "ports", "", "Comma-separated list of specific ports to scan")
}
func main() {
        flag.Parse()



	
}