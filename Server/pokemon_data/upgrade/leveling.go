package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"

	models "Net-Centric--Pokemon/Server/pokemon_data/models"
)

// UpdateStats adjusts Pokémon stats after leveling up.
func UpdateStats(p *models.Pokemon) {
	p.Stats.HP = int(math.Round(float64(p.Stats.HP) * (1 + p.EV)))
	p.Stats.Attack = int(math.Round(float64(p.Stats.Attack) * (1 + p.EV)))
	p.Stats.Defense = int(math.Round(float64(p.Stats.Defense) * (1 + p.EV)))
	p.Stats.SpAtk = int(math.Round(float64(p.Stats.SpAtk) * (1 + p.EV)))
	p.Stats.SpDef = int(math.Round(float64(p.Stats.SpDef) * (1 + p.EV)))
	p.Stats.Total = p.Stats.HP + p.Stats.Attack + p.Stats.Defense + p.Stats.SpAtk + p.Stats.SpDef + p.Stats.Speed
}

// GainExp awards experience points to a Pokémon. Thread-safe.
func GainExp(p *models.Pokemon, exp int, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	p.Stats.Exp += exp
	fmt.Printf("Pokemon %s gained %d exp\n", p.Name, exp)
}

// LevelUp levels up a Pokémon if enough experience is accumulated. Thread-safe.
func LevelUp(p *models.Pokemon, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	neededExp := int(math.Pow(2, float64(p.Level)) * float64(p.Stats.Total))
	if p.Stats.Exp >= neededExp {
		p.Level++
		UpdateStats(p)
		fmt.Printf("Pokemon %s leveled up to %d!\n", p.Name, p.Level)
		return nil
	}
	return fmt.Errorf("not enough exp to level up, need %d, have %d", neededExp, p.Stats.Exp)
}

// LoadPokemons loads Pokémon data from a JSON file.
func LoadPokemons(filePath string) ([]models.Pokemon, map[string]models.Pokemon, error) {
	pokemonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read Pokémon data: %w", err)
	}

	var pokemons []models.Pokemon
	err = json.Unmarshal(pokemonData, &pokemons)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal Pokémon data: %w", err)
	}

	pokemonByNumber := make(map[string]models.Pokemon)
	for _, pokemon := range pokemons {
		pokemonByNumber[pokemon.Number] = pokemon
	}

	return pokemons, pokemonByNumber, nil
}

// MainLeveling demonstrates leveling and evolution functionalities, supporting concurrency.
func MainLeveling(wg *sync.WaitGroup, evolutionChains map[string][]models.Pokemon) {
	defer wg.Done()

	// Load Pokémon data
	pokemons, _, err := LoadPokemons("./pokemon_data/pokedex_data/pokemon_types.json")
	if err != nil {
		fmt.Printf("Error loading Pokémon data: %v\n", err)
		return
	}

	mu := &sync.Mutex{}

	// Example Pokémon for leveling and evolution
	examplePokemon := pokemons[0]

	// Simulate gaining experience and leveling up
	go GainExp(&examplePokemon, 500, mu)
	go func() {
		err := LevelUp(&examplePokemon, mu)
		if err != nil {
			fmt.Printf("LevelUp error: %v\n", err)
		}
	}()

	// Simulate evolution
	go func() {
		err := Evolve(&examplePokemon, evolutionChains)
		if err != nil {
			fmt.Printf("Evolution error: %v\n", err)
		} else {
			fmt.Printf("Evolved Pokemon: %+v\n", examplePokemon)
		}
	}()
}
