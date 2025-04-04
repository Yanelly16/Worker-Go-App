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

func scanTarget(target string, portList []int) {
    startTime := time.Now()
    tasks := make(chan string, workers)
    results := make(chan ScanResult)
    var openPorts []ScanResult
    var wg sync.WaitGroup

    dialer := net.Dialer{Timeout: time.Duration(timeout) * time.Second}

    // Worker pool setup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go worker(&wg, tasks, results, dialer)
    }

    // Task distribution
    go func() {
        if len(portList) > 0 {
            for _, port := range portList {
                tasks <- net.JoinHostPort(target, strconv.Itoa(port))
            }
        } else {
            for port := startPort; port <= endPort; port++ {
                tasks <- net.JoinHostPort(target, strconv.Itoa(port))
            }
        }
        close(tasks)
    }()

    // Result collection
    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        if result.State == "open" {
            openPorts = append(openPorts, result)
        }
    }
}

func worker(wg *sync.WaitGroup, tasks <-chan string, results chan<- ScanResult, dialer net.Dialer) {
    defer wg.Done()
    for addr := range tasks {
        port, _ := strconv.Atoi(strings.Split(addr, ":")[1])
        result := ScanResult{Port: port}

        conn, err := dialer.Dial("tcp", addr)
        if err != nil {
            result.State = "closed"
        } else {
            result.State = "open"
            buffer := make([]byte, 1024)
            conn.SetReadDeadline(time.Now().Add(dialer.Timeout))
            n, _ := conn.Read(buffer)
            if n > 0 {
                result.Banner = strings.TrimSpace(string(buffer[:n]))
            }
            conn.Close()
        }
        results <- result
    }
}
func printTextSummary(summary ScanSummary) {
    fmt.Println("\nScan Summary:")
    fmt.Printf("Target: %s\n", summary.Target)
    if summary.StartPort != 0 && summary.EndPort != 0 {
        fmt.Printf("Port range: %d-%d\n", summary.StartPort, summary.EndPort)
    }
    fmt.Printf("Ports scanned: %d\n", summary.PortsScanned)
    fmt.Printf("Open ports: %d\n", summary.OpenPorts)
    fmt.Printf("Time taken: %v\n", summary.TimeTaken)

    if len(summary.Results) > 0 {
        fmt.Println("\nOpen Ports:")
        for _, result := range summary.Results {
            fmt.Printf("  %d: %s\n", result.Port, result.Banner)
        }
    }
}
func main() {
        flag.Parse()
        targetList := strings.Split(targets, ",")
    var portList []int

    if portsList != "" {
        for _, p := range strings.Split(portsList, ",") {
            port, _ := strconv.Atoi(strings.TrimSpace(p))
            portList = append(portList, port)
        }
    }

    for _, target := range targetList {
        startTime := time.Now()
        openPorts := scanTarget(target, portList)
        summary := ScanSummary{
            Target:       target,
            StartPort:    startPort,
            EndPort:      endPort,
            PortsScanned: len(portList),
            OpenPorts:    len(openPorts),
            TimeTaken:    time.Since(startTime).String(),
            Results:      openPorts,
        }

        if jsonOutput {
            jsonData, _ := json.MarshalIndent(summary, "", "  ")
            fmt.Println(string(jsonData))
        } else {
            printTextSummary(summary)
        }
    }


	
}