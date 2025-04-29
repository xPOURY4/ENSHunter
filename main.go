package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

const SimpleRegistrarControllerABI = `[
    {
        "constant": true,
        "inputs": [
            {
                "name": "name",
                "type": "string"
            }
        ],
        "name": "available",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }
]`

type RegistrarController struct {
	*bind.BoundContract
}

func NewRegistrarController(address common.Address, backend bind.ContractBackend) (*RegistrarController, error) {
	parsedABI, err := abi.JSON(strings.NewReader(SimpleRegistrarControllerABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, parsedABI, backend, backend, backend)
	return &RegistrarController{contract}, nil
}

func (rc *RegistrarController) Available(opts *bind.CallOpts, name string) (bool, error) {
	var out []interface{}
	err := rc.BoundContract.Call(opts, &out, "available", name)
	if err != nil {
		return false, err
	}
	return out[0].(bool), nil
}

type Config struct {
	InfuraKey string `json:"infura_key"`
	Workers   int    `json:"workers"`
	RateLimit int    `json:"rate_limit"`
	Retries   int    `json:"retries"`
	Timeout   int    `json:"timeout"`
}

func loadConfig() Config {
	config := Config{
		Workers:   5,
		RateLimit: 100,
		Retries:   3,
		Timeout:   30,
	}

	configDir := getConfigDir()
	configFile := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configFile); err == nil {
		data, err := ioutil.ReadFile(configFile)
		if err == nil {
			json.Unmarshal(data, &config)
		}
	}

	if err := godotenv.Load(); err == nil {
		if key := os.Getenv("INFURA_KEY"); key != "" {
			config.InfuraKey = key
		}
		if workers := os.Getenv("WORKERS"); workers != "" {
			fmt.Sscanf(workers, "%d", &config.Workers)
		}
		if rate := os.Getenv("RATE_LIMIT"); rate != "" {
			fmt.Sscanf(rate, "%d", &config.RateLimit)
		}
		if retries := os.Getenv("RETRIES"); retries != "" {
			fmt.Sscanf(retries, "%d", &config.Retries)
		}
		if timeout := os.Getenv("TIMEOUT"); timeout != "" {
			fmt.Sscanf(timeout, "%d", &config.Timeout)
		}
	}

	return config
}

func saveConfig(config Config) error {
	configDir := getConfigDir()
	
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}
	
	configFile := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(configFile, data, 0644)
}

func getConfigDir() string {
	usr, err := user.Current()
	if err != nil {
		return ".enshunter"
	}
	return filepath.Join(usr.HomeDir, ".enshunter")
}

type Result struct {
	Domain      string
	IsAvailable bool
	Error       error
}

func main() {
	// Load configuration from files
	config := loadConfig()
	
	// Initialize color objects
	successColor := color.New(color.FgGreen, color.Bold)
	errorColor := color.New(color.FgRed, color.Bold)
	infoColor := color.New(color.FgCyan)
	
	// Use the color objects directly to avoid "imported and not used" error
	successColor.Print("ENSHunter ")
	fmt.Println("initialized")
	
	// Command line flags
	infuraKey := flag.String("infura", config.InfuraKey, "Infura Project ID")
	inputFile := flag.String("input", "esn.txt", "Input file containing domain names")
	outputFile := flag.String("output", "ens_available.txt", "Output file for available domains")
	workers := flag.Int("workers", config.Workers, "Number of concurrent workers")
	rateLimit := flag.Int("rate", config.RateLimit, "Rate limit in milliseconds between requests")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging")
	retries := flag.Int("retries", config.Retries, "Number of retries for failed requests")
	timeout := flag.Int("timeout", config.Timeout, "Request timeout in seconds")
	saveConfigFlag := flag.Bool("save-config", false, "Save current settings as default configuration")
	flag.Parse()

	// Save configuration if requested
	if *saveConfigFlag {
		newConfig := Config{
			InfuraKey: *infuraKey,
			Workers:   *workers,
			RateLimit: *rateLimit,
			Retries:   *retries,
			Timeout:   *timeout,
		}
		
		if err := saveConfig(newConfig); err != nil {
			errorColor.Println("Warning: Failed to save configuration:", err)
		} else {
			successColor.Println("Configuration saved successfully!")
		}
	}

	if *infuraKey == "" {
		errorColor.Println("Infura Project ID is required. Use -infura flag or set INFURA_KEY in .env file or config.json")
		os.Exit(1)
	}

	ethereumNode := fmt.Sprintf("https://mainnet.infura.io/v3/%s", *infuraKey)
	verbose := *verboseFlag
	
	if verbose {
		infoColor.Println("Connecting to Ethereum network...")
	}
	
	client, err := ethclient.Dial(ethereumNode)
	if err != nil {
		errorColor.Println("Failed to connect to Ethereum:", err)
		os.Exit(1)
	}

	if verbose {
		successColor.Println("Successfully connected to Ethereum")
	}

	registrarAddress := common.HexToAddress("0x283Af0B28c62C092C9727F1Ee09c02CA627EB7F5")
	controller, err := NewRegistrarController(registrarAddress, client)
	if err != nil {
		errorColor.Println("Failed to create controller:", err)
		os.Exit(1)
	}

	domains, err := loadDomains(*inputFile)
	if err != nil {
		errorColor.Println("Failed to load domains:", err)
		os.Exit(1)
	}

	if len(domains) == 0 {
		errorColor.Println("No domains found in input file")
		os.Exit(1)
	}

	outFile, err := os.Create(*outputFile)
	if err != nil {
		errorColor.Println("Failed to create output file:", err)
		os.Exit(1)
	}
	defer outFile.Close()
	
	outputLock := &sync.Mutex{}
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	var (
		wg sync.WaitGroup
		available int
		errorCount int
	)

	jobs := make(chan string, len(domains))
	results := make(chan Result, len(domains))
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	fmt.Print("Starting ENSHunter - checking ")
	infoColor.Printf("%d", len(domains))
	fmt.Println(" domains")
	
	bar := progressbar.Default(int64(len(domains)))

	for w := 0; w < *workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range jobs {
				processedDomain := strings.TrimSuffix(domain, ".eth")
				
				var isAvailable bool
				var err error
				
				for attempt := 0; attempt <= *retries; attempt++ {
					if attempt > 0 && verbose {
						infoColor.Printf("Retry %d for domain %s\n", attempt, domain)
					}
					
					isAvailable, err = controller.Available(&bind.CallOpts{Context: ctx}, processedDomain)
					if err == nil {
						break
					}
					
					time.Sleep(time.Duration(*rateLimit) * time.Millisecond * 2)
				}
				
				results <- Result{
					Domain:      domain,
					IsAvailable: isAvailable,
					Error:       err,
				}
				
				time.Sleep(time.Duration(*rateLimit) * time.Millisecond)
			}
		}()
	}

	for _, domain := range domains {
		jobs <- domain
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		bar.Add(1)
		
		if result.Error != nil {
			if verbose {
				errorColor.Printf("Error checking %s: %v\n", result.Domain, result.Error)
			}
			errorCount++
			continue
		}

		if result.IsAvailable {
			available++
			if verbose {
				successColor.Printf("Domain %s is available\n", result.Domain)
			}
			
			outputLock.Lock()
			_, err := writer.WriteString(result.Domain + "\n")
			if err != nil && verbose {
				errorColor.Printf("Error writing to file: %v\n", err)
			}
			writer.Flush()
			outputLock.Unlock()
		} else if verbose {
			fmt.Printf("Domain %s is not available\n", result.Domain)
		}
	}

	fmt.Println()
	successColor.Println("Scan completed!")
	fmt.Print("Total domains checked: ")
	infoColor.Printf("%d\n", len(domains))
	fmt.Print("Available domains: ")
	successColor.Printf("%d\n", available)
	fmt.Print("Errors: ")
	if errorCount > 0 {
		errorColor.Printf("%d\n", errorCount)
	} else {
		successColor.Printf("%d\n", errorCount)
	}
	fmt.Print("Available domains saved to: ")
	infoColor.Println(*outputFile)
}

func loadDomains(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}
		
		if !strings.HasSuffix(domain, ".eth") {
			domain = domain + ".eth"
		}
		
		domains = append(domains, domain)
	}
	
	return domains, scanner.Err()
}
