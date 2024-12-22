package main

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/gocolly/colly/v2"
)

type Pokemon struct {
	Name  string   `json:"name"`
	Types []string `json:"types"`
	Stats Stats    `json:"stats"`
}

type Stats struct {
	Total   int `json:"total"`
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
	pokemonByType := make(map[string][]Pokemon)

	c := colly.NewCollector()

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		name := e.ChildText("a.ent-name")
		if name == "" {
			return
		}

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
			HP:      toInt(hp),
			Attack:  toInt(attack),
			Defense: toInt(defense),
			SpAtk:   toInt(spAtk),
			SpDef:   toInt(spDef),
			Speed:   toInt(speed),
		}

		for _, pokemonType := range types {
			if len(pokemonByType[pokemonType]) < typeLimit {
				pokemonByType[pokemonType] = append(pokemonByType[pokemonType], Pokemon{
					Name:  name,
					Types: types,
					Stats: stats,
				})
			}
		}


	})

	fmt.Println("Collecting Pokemon data...")
	c.Visit("https://pokemondb.net/pokedex/all")

	limitedPokemons := []Pokemon{}
	for _, pokemons := range pokemonByType {
		limitedPokemons = append(limitedPokemons, pokemons...)
	}

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
}

func toInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}
