// Filename: main.go
// Purpose: Concurrent TCP port scanner with configurable options and enhanced features
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"

	//"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ScanResult represents the result of scanning a single port
type ScanResult struct {
	Port   int    `json:"port"`             // Port number that was scanned
	State  string `json:"state"`            // State of the port ("open" or "closed")
	Banner string `json:"banner,omitempty"` // Banner/service information if port is open (optional in JSON)
}

// ScanSummary aggregates all scan results for a target
type ScanSummary struct {
	Target       string       `json:"target"`               // The host that was scanned
	StartPort    int          `json:"start_port,omitempty"` // First port in range (optional in JSON)
	EndPort      int          `json:"end_port,omitempty"`   // Last port in range (optional in JSON)
	PortsScanned int          `json:"ports_scanned"`        // Total number of ports scanned
	OpenPorts    int          `json:"open_ports"`           // Number of open ports found
	TimeTaken    string       `json:"time_taken"`           // Duration of the scan
	Results      []ScanResult `json:"results,omitempty"`    // Detailed scan results (optional in JSON)
}

// Command-line flags
var (
	targets    string // Comma-separated list of targets to scan
	startPort  int    // First port in range to scan
	endPort    int    // Last port in range to scan
	workers    int    // Number of concurrent workers
	timeout    int    // Connection timeout in seconds
	jsonOutput bool   // Flag for JSON output format
	portsList  string // Comma-separated list of specific ports to scan
)

// init initializes command-line flags
func init() {
	flag.StringVar(&targets, "targets", "scanme.nmap.org", "Comma-separated list of targets")
	flag.IntVar(&startPort, "start-port", 1, "Start port number")
	flag.IntVar(&endPort, "end-port", 1024, "End port number")
	flag.IntVar(&workers, "workers", 100, "Number of concurrent workers")
	flag.IntVar(&timeout, "timeout", 5, "Connection timeout in seconds")
	flag.BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	flag.StringVar(&portsList, "ports", "", "Comma-separated list of specific ports to scan")
}

// main is the entry point of the program
func main() {
	flag.Parse() // Parse command-line flags

	// Split comma-separated targets into a slice
	targetList := strings.Split(targets, ",")
	var portList []int // Will hold specific ports if provided

	// Process specific ports list if provided
	if portsList != "" {
		for _, p := range strings.Split(portsList, ",") {
			port, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				fmt.Printf("Invalid port number: %s\n", p)
				continue
			}
			portList = append(portList, port)
		}
	}

	// Scan each target in the list
	for _, target := range targetList {
		scanTarget(target, portList)
	}
}

// scanTarget handles scanning a single target host
func scanTarget(target string, portList []int) {
	startTime := time.Now() // Record start time for duration calculation

	// Create channels for task distribution and result collection
	tasks := make(chan string, workers) // Buffered channel for addresses to scan
	results := make(chan ScanResult)    // Channel for scan results
	var openPorts []ScanResult          // Slice to store open port results
	var wg sync.WaitGroup               // WaitGroup to track worker completion

	// Configure dialer with timeout
	dialer := net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Create worker pool
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, results, dialer)
	}

	// Send ports to scan (runs in separate goroutine)
	go func() {
		if len(portList) > 0 {
			// Scan specific ports if provided
			for _, port := range portList {
				address := net.JoinHostPort(target, strconv.Itoa(port))
				tasks <- address
			}
		} else {
			// Scan port range otherwise
			for port := startPort; port <= endPort; port++ {
				address := net.JoinHostPort(target, strconv.Itoa(port))
				tasks <- address
				// Show progress every 100 ports (for text output)
				if port%100 == 0 && !jsonOutput {
					fmt.Printf("\rScanning port %d/%d...", port, endPort)
				}
			}
		}
		close(tasks) // Close tasks channel when done
	}()

	// Wait for workers to finish and close results channel (runs in separate goroutine)
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results from workers
	for result := range results {
		if result.State == "open" {
			openPorts = append(openPorts, result)
			// Print open ports immediately for text output
			if !jsonOutput {
				fmt.Printf("\rPort %d is %s - %s\n", result.Port, result.State, result.Banner)
			}
		}
	}

	// Prepare summary statistics
	duration := time.Since(startTime)
	portsScanned := endPort - startPort + 1
	if len(portList) > 0 {
		portsScanned = len(portList)
	}

	// Create summary struct
	summary := ScanSummary{
		Target:       target,
		StartPort:    startPort,
		EndPort:      endPort,
		PortsScanned: portsScanned,
		OpenPorts:    len(openPorts),
		TimeTaken:    duration.String(),
		Results:      openPorts,
	}

	// Output results in requested format
	if jsonOutput {
		jsonData, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			fmt.Println("Error generating JSON output:", err)
			return
		}
		fmt.Println(string(jsonData))
	} else {
		printTextSummary(summary)
	}
}

// worker is a goroutine that performs the actual port scanning
func worker(wg *sync.WaitGroup, tasks <-chan string, results chan<- ScanResult, dialer net.Dialer) {
	defer wg.Done() // Signal completion when done

	// Process each address from the tasks channel
	for addr := range tasks {
		// Extract port from address (e.g., "example.com:80" -> 80)
		port, _ := strconv.Atoi(strings.Split(addr, ":")[1])
		result := ScanResult{Port: port}

		// Attempt TCP connection
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			result.State = "closed"
		} else {
			result.State = "open"
			// Try to read banner if port is open
			buffer := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(dialer.Timeout))
			n, _ := conn.Read(buffer)
			if n > 0 {
				result.Banner = strings.TrimSpace(string(buffer[:n]))
			}
			conn.Close()
		}
		// Send result back through channel
		results <- result
	}
}

// printTextSummary outputs results in human-readable format
func printTextSummary(summary ScanSummary) {
	fmt.Println("\nScan Summary:")
	fmt.Printf("Target: %s\n", summary.Target)
	if summary.StartPort != 0 && summary.EndPort != 0 {
		fmt.Printf("Port range: %d-%d\n", summary.StartPort, summary.EndPort)
	}
	fmt.Printf("Ports scanned: %d\n", summary.PortsScanned)
	fmt.Printf("Open ports: %d\n", summary.OpenPorts)
	fmt.Printf("Time taken: %v\n", summary.TimeTaken)

	// Print detailed results if any open ports found
	if len(summary.Results) > 0 {
		fmt.Println("\nOpen Ports:")
		for _, result := range summary.Results {
			fmt.Printf("  %d: %s\n", result.Port, result.Banner)
		}
	}
}