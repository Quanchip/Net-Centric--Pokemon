package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/gocolly/colly/v2"
)

type Pokemon struct {
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	Number  string   `json:"number"`
	SubName string   `json:"sub_name"`
	Stats   Stats    `json:"stats"`
}

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

func main() {
	// Map to store Pokemon by type, with a limit of 5 per type
	typeLimit := 5
	pokemonByNumber := make(map[string]Pokemon)
	pokemonByType := make(map[string][]Pokemon)

	c := colly.NewCollector()

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		name := e.ChildText("a.ent-name")
		if name == "" {
			return
		}
		subName := e.ChildText("small.text-muted")

		// Get types for this Pokemon
		types := []string{}
		e.ForEach("td.cell-icon a.type-icon", func(_ int, el *colly.HTMLElement) {
			pokemonType := el.Text
			if pokemonType != "" {
				// // Only add if we haven't reached the limit for this type
				// if len(pokemonByType[pokemonType]) < typeLimit {
				// 	pokemonByType[pokemonType] = append(pokemonByType[pokemonType], name)
				// 	pokemonByType[pokemonType] = append(pokemonByType[pokemonType], stats)
				// }
				types = append(types, pokemonType)
			}
		})

		// Extract number
		number := e.ChildText("td.cell-num span.infocard-cell-data")

		// Extract stats

		total := e.ChildText("td.cell-num.cell-total")
		hp := e.ChildText("td:nth-of-type(5)")
		attack := e.ChildText("td:nth-of-type(6)")
		defense := e.ChildText("td:nth-of-type(7)")
		spAtk := e.ChildText("td:nth-of-type(8)")
		spDef := e.ChildText("td:nth-of-type(9)")
		speed := e.ChildText("td:nth-of-type(10)")

		// Convert stats to integers
		stats := Stats{
			Total:   toInt(total),
			Exp:     toInt(total), // Assign Total as Exp in Stats
			HP:      toInt(hp),
			Attack:  toInt(attack),
			Defense: toInt(defense),
			SpAtk:   toInt(spAtk),
			SpDef:   toInt(spDef),
			Speed:   toInt(speed),
		}

		pokemon := Pokemon{
			Name:    name,
			SubName: subName,
			Types:   types,
			Number:  number,
			Stats:   stats,
		}
		pokemonByNumber[number] = pokemon

		for _, pokemonType := range types {
			// skip if a pokemon of the same type has been saved already
			typeAlreadySaved := false
			for _, savedPokemon := range pokemonByType[pokemonType] {
				if arePokemonsEqual(savedPokemon, pokemon) {
					typeAlreadySaved = true
					break
				}
			}

			if !typeAlreadySaved {
				if len(pokemonByType[pokemonType]) < typeLimit {
					pokemonByType[pokemonType] = append(pokemonByType[pokemonType], pokemon)
				}
			}
		}

	})

	fmt.Println("Collecting Pokemon data...")
	c.Visit("https://pokemondb.net/pokedex/all")

	limitedPokemons := []Pokemon{}
	for _, pokemons := range pokemonByType {
		limitedPokemons = append(limitedPokemons, pokemons...)
	}

	// Sort by number, will change the current array
	sort.Slice(limitedPokemons, func(i, j int) bool {
		numI, _ := strconv.Atoi(limitedPokemons[i].Number)
		numJ, _ := strconv.Atoi(limitedPokemons[j].Number)
		return numI < numJ
	})

	// Save to JSON file
	file, _ := os.Create("./pokedex_data/pokemon_types.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(limitedPokemons)

	fmt.Println("Data saved to pokemon_data.json")

	// Print the results
	// fmt.Println("\nCollected Pokemon by type (max 5 per type):")
	// for typeName, pokemons := range pokemonByType {
	//     fmt.Printf("\n%s type (%d Pokemon):\n", typeName, len(pokemons))
	//     for _, name := range pokemons {
	//         fmt.Printf("- %s\n", name)
	//     }
	// }

	evolutionChains := findEvolutionChains(limitedPokemons, pokemonByNumber, 3)

	// Save to JSON file
	fileEvolution, _ := os.Create("./pokedex_data/pokemon_evolution.json")
	defer fileEvolution.Close()

	encoderEvolution := json.NewEncoder(fileEvolution)
	encoderEvolution.SetIndent("", "  ")
	encoderEvolution.Encode(evolutionChains)

	fmt.Println("Evolution data saved to pokemon_evolution.json")
}

func toInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

// Method to find the pokemon that near by the other number under the range
func findEvolutionChains(pokemons []Pokemon, pokemonByNumber map[string]Pokemon, rangeVal int) map[string][]Pokemon {
	evolutionChains := make(map[string][]Pokemon)

	for _, pokemon := range pokemons {
		chain := []Pokemon{pokemon}
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

func arePokemonsEqual(p1, p2 Pokemon) bool {
	p1JSON, _ := json.Marshal(p1)
	p2JSON, _ := json.Marshal(p2)
	return string(p1JSON) == string(p2JSON)
}
