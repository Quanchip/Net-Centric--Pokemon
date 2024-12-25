package main

import (
	"Net-Centric--Pokemon/Server/pokemon_data/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"sync" // Import the sync package
	"time"
)

// Pokemon Data (Moved here)
var allPokemons []models.Pokemon

// Player data (Moved Here)
type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerRequest struct {
	Name string `json:"name"`
}

var (
	players       = make(map[int]*Player)
	playersMutex  sync.Mutex
	playerCounter = 0
)

func loadPokemonData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var pokemons []models.Pokemon
	if err := decoder.Decode(&pokemons); err != nil {
		return err
	}

	allPokemons = pokemons
	return nil

}

// Randomly select 5 Pokémon
func getRandomPokemons() []models.Pokemon {
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(allPokemons))
	selected := []models.Pokemon{}

	for i := 0; i < len(perm) && len(selected) < 5; i++ {
		pokemon := allPokemons[perm[i]]
		alreadySelected := false
		for _, p := range selected {
			if p.Name == pokemon.Name {
				alreadySelected = true
				break
			}
		}
		if !alreadySelected {
			selected = append(selected, pokemon)
		}
	}
	return selected
}

// handle connection when player connect
func handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Đọc request body
	var req PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Tạo player mới
	playersMutex.Lock()
	playerCounter++
	player := &Player{
		ID:   playerCounter,
		Name: req.Name,
	}
	players[player.ID] = player
	playersMutex.Unlock()

	fmt.Printf("New player connected. ID: %d, Name: %s\n", player.ID, player.Name)

	// Trả về thông tin player
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// Handler for /join endpoint
func joinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	selectedPokemons := getRandomPokemons()
	response := struct {
		Pokemons []models.Pokemon `json:"pokemons"`
	}{
		Pokemons: selectedPokemons,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Use WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	// Initialize the pokemon data
	if err := loadPokemonData("./pokemon_data/pokedex_data/pokemon_types.json"); err != nil {
		fmt.Println("Error loading Pokémon data:", err)
		return
	}
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
		if err := loadPokemonData("./pokemon_data/pokedex_data/pokemon_types.json"); err != nil {
			fmt.Println("Error reloading Pokémon data:", err)
			return
		}
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

	// Set up HTTP server
	http.HandleFunc("/connect", handleConnect)
	http.HandleFunc("/join", joinHandler)

	// Wait for the processes to complete
	wg.Wait()
	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}

}
