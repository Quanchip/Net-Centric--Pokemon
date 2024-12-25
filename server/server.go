package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync" // Import the sync package
	"time"
)

func main() {
	// Use WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Run scraper to update `pokemon_types.json`
	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		fmt.Println("Starting scraper...")

		cmd := exec.CommandContext(ctx, "go", "run", "./pokemon_data/pokedex.go")
		// Print command output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("Error creating stdout pipe %v\n", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Printf("Error creating stderr pipe %v\n", err)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Printf("Error running go file: %v\n", err)
			return
		}
		// Combine stderr, and stdout so it is printed to the same console
		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stderr, stderr)
		if err := cmd.Wait(); err != nil {
			log.Printf("Error waiting for command completion: %v\n", err)
		}

		fmt.Println("Scraper finished.")
	}()

	// After scraping, calculate levels and save them
	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		fmt.Println("Starting leveling...")
		cmd := exec.CommandContext(ctx, "go", "run", "./pokemon_data/upgrade/leveling.go")
		// Print command output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("Error creating stdout pipe %v\n", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Printf("Error creating stderr pipe %v\n", err)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Printf("Error running go file: %v\n", err)
			return
		}
		// Combine stderr, and stdout so it is printed to the same console
		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stderr, stderr)

		if err := cmd.Wait(); err != nil {
			log.Printf("Error waiting for command completion: %v\n", err)
		}
		fmt.Println("Leveling finished.")
	}()

	// Wait for the processes to complete
	wg.Wait()

}
