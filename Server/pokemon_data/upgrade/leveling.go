package main

import (
	// Import the models package
	models "Net-Centric--Pokemon/Server/pokemon_data/models"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func findEvolutionChains(pokemons []models.Pokemon, pokemonByNumber map[string]models.Pokemon, rangeVal int) map[string][]models.Pokemon {
	evolutionChains := make(map[string][]models.Pokemon)

	for _, pokemon := range pokemons {
		chain := []models.Pokemon{pokemon}
		num, err := strconv.Atoi(pokemon.Number)
		if err != nil {
			fmt.Println("Cannot parse string to int")
		}

		for i := 1; i <= rangeVal; i++ {
			if len(chain) >= 3 {
				break
			}
			nextNumStr := fmt.Sprintf("%04d", num+i)
			nextPokemon, ok := pokemonByNumber[nextNumStr]

			// We only add pokemons which have higher number, and are the same type
			if ok && areSameTypes(pokemon.Types, nextPokemon.Types) {
				chain = append(chain, nextPokemon)
			}
		}
		if len(chain) > 1 {
			evolutionChains[pokemon.Name] = chain
		}
	}
	return evolutionChains
}

func areSameTypes(types1, types2 []string) bool {
	if len(types1) == 0 || len(types2) == 0 {
		return false
	}

	// Single type: match any of types in types2
	if len(types1) == 1 {
		for _, typ1 := range types1 {
			for _, typ2 := range types2 {
				if typ1 == typ2 {
					return true
				}
			}
		}
		return false
	}

	// Double type: match all types in types2
	if len(types1) == 2 {
		typeMap := make(map[string]bool)
		for _, typ1 := range types1 {
			typeMap[typ1] = true
		}

		for _, typ2 := range types2 {
			if !typeMap[typ2] {
				return false
			}
		}
		return true
	}

	// Previous pokemon single, and current has two
	if len(types1) == 1 && len(types2) == 2 {
		for _, typ1 := range types1 {
			typeMatch := false
			for _, typ2 := range types2 {
				if typ1 == typ2 {
					typeMatch = true
				}

			}
			if !typeMatch {
				return false
			}
		}
		return true
	}

	return false

}

func main() {
	mainLeveling()
}

func mainLeveling() {
	// Load pokemon from json file
	pokemonData, err := os.ReadFile("./pokemon_data/pokedex_data/pokemon_types.json")
	if err != nil {
		fmt.Println("Failed to read pokemon data:", err)
		return
	}

	var limitedPokemons []models.Pokemon
	err = json.Unmarshal(pokemonData, &limitedPokemons)
	if err != nil {
		fmt.Println("Failed to unmarshal pokemon data:", err)
		return
	}

	pokemonByNumber := make(map[string]models.Pokemon)

	for _, pokemon := range limitedPokemons {
		pokemonByNumber[pokemon.Number] = pokemon
	}

	evolutionChains := findEvolutionChains(limitedPokemons, pokemonByNumber, 3)

	// Save to JSON file
	fileEvolution, err := os.Create("./pokemon_data/pokedex_data/pokemon_evolution.json")
	if err != nil {
		return
	}
	defer fileEvolution.Close()

	encoderEvolution := json.NewEncoder(fileEvolution)
	encoderEvolution.SetIndent("", "  ")
	encoderEvolution.Encode(evolutionChains)

	fmt.Println("Evolution data saved to pokemon_evolution.json")
}
