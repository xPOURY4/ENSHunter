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
	"sync"
	"strings"
	"time"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
	"github.com/fatih/color"
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
	abi, err := abi.JSON(strings.NewReader(SimpleRegistrarControllerABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, abi, backend, backend, backend)
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
	Domain     string
	IsAvailable bool
	Error      error
}

func main() {
	// Load configuration from files
	config := loadConfig()
	
	// Create styled output with the fatih/color package
	successColor := color.New(color.FgGreen, color.Bold)
	errorColor := color.New(color.FgRed, color.Bold)
	infoColor := color.New(color.FgCyan)
	
	// Use the color objects directly to avoid "imported and not used" error
	successColor.Println("ENSHunter initialized")
	
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
		config := Config{
			InfuraKey: *infuraKey,
			Workers:   *workers,
			RateLimit: *rateLimit,
			Retries:   *retries,
			Timeout:   *timeout,
		}
		
		if err := saveConfig(config); err != nil {
			log.Printf(errorColor.Sprintf("Warning: Failed to save configuration: %v", err))
		} else {
			log.Println(successColor.Sprintf("Configuration saved successfully!"))
		}
	}

	if *infuraKey == "" {
		log.Fatal(errorColor.Sprintf("Infura Project ID is required. Use -infura flag or set INFURA_KEY in .env file or config.json"))
	}

	ethereumNode := fmt.Sprintf("https://mainnet.infura.io/v3/%s", *infuraKey)
	verbose := *verboseFlag
	
	if verbose {
		log.Println(infoColor.Sprintf("Connecting to Ethereum network..."))
	}
	
	client, err := ethclient.Dial(ethereumNode)
	if err != nil {
		log.Fatalf(errorColor.Sprintf("Failed to connect to Ethereum: %v", err))
	}

	if verbose {
		log.Println(successColor.Sprintf("Successfully connected to Ethereum"))
	}

	registrarAddress := common.HexToAddress("0x283Af0B28c62C092C9727F1Ee09c02CA627EB7F5")
	controller, err := NewRegistrarController(registrarAddress, client)
	if err != nil {
		log.Fatalf(errorMsg("Failed to create controller: %v"), err)
	}

	domains, err := loadDomains(*inputFile)
	if err != nil {
		log.Fatalf(errorMsg("Failed to load domains: %v"), err)
	}

	if len(domains) == 0 {
		log.Fatal(errorMsg("No domains found in input file"))
	}

	outFile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf(errorMsg("Failed to create output file: %v"), err)
	}
	defer outFile.Close()
	
	outputLock := &sync.Mutex{}
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	var (
		wg sync.WaitGroup
		available int32
		errorCount int32
	)

	jobs := make(chan string, len(domains))
	results := make(chan Result, len(domains))
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	fmt.Printf(infoColor.Sprintf("Starting ENSHunter - checking %d domains\n"), len(domains))
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
						log.Printf(infoColor.Sprintf("Retry %d for domain %s"), attempt, domain)
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
				log.Printf(errorColor.Sprintf("Error checking %s: %v"), result.Domain, result.Error)
			}
			atomic.AddInt32(&errorCount, 1)
			continue
		}

		if result.IsAvailable {
			atomic.AddInt32(&available, 1)
			if verbose {
				log.Printf(successColor.Sprintf("Domain %s is available"), result.Domain)
			}
			
			outputLock.Lock()
			_, err := writer.WriteString(result.Domain + "\n")
			if err != nil && verbose {
				log.Printf(errorColor.Sprintf("Error writing to file: %v"), err)
			}
			writer.Flush()
			outputLock.Unlock()
		} else if verbose {
			log.Printf("Domain %s is not available", result.Domain)
		}
	}

	fmt.Printf("\n%s\n", successColor.Sprintf("Scan completed!"))
	fmt.Printf("Total domains checked: %s\n", infoColor.Sprintf("%d", len(domains)))
	fmt.Printf("Available domains: %s\n", successColor.Sprintf("%d", available))
	fmt.Printf("Errors: %s\n", errorCount > 0 ? errorColor.Sprintf("%d", errorCount) : successColor.Sprintf("%d", errorCount))
	fmt.Printf("Available domains saved to: %s\n", infoColor.Sprintf(*outputFile))
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
