// SERVER CODE
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type Stats struct {
	Total   int `json:"total"`
	Exp     int `json:"exp"`
	HP      int `json:"hp"`
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	SpAtk   int `json:"sp_atk"`
	SpDef   int `json:"sp_def"`
	Speed   int `json:"speed"`
}

type Pokemon struct {
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	Number  string   `json:"number"`
	SubName string   `json:"sub_name"`
	Stats   Stats    `json:"stats"`
}

type EvolutionChain struct {
	Chain []Pokemon `json:"chain"`
}

var allPokemons []Pokemon

type Player struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type PlayerRequest struct {
    Name string `json:"name"`
}

var (
    players      = make(map[int]*Player)
    playersMutex sync.Mutex
    playerCounter = 0
)

// Load Pokémon data from JSON file
func loadPokemonData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	data := make(map[string][]Pokemon)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	for _, chain := range data {
		allPokemons = append(allPokemons, chain...)
	}
	return nil
}

// Randomly select 5 Pokémon
func getRandomPokemons() []Pokemon {
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(allPokemons))
	selected := []Pokemon{}

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
		Pokemons []Pokemon `json:"pokemons"`
	}{
		Pokemons: selectedPokemons,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Load Pokémon data
	if err := loadPokemonData(".././pokedex_data/pokemon_evolution.json"); err != nil {
		fmt.Println("Error loading Pokémon data:", err)
		return
	}

	// Set up HTTP server
	http.HandleFunc("/connect", handleConnect)
	http.HandleFunc("/join", joinHandler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
