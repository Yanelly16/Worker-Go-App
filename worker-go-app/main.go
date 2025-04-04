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

func main() {




	
}