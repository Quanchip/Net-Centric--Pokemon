package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	models "Net-Centric--Pokemon/Server/pokemon_data/models"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Map to store Pokemon by type, with a limit of 5 per type
	typeLimit := 5
	pokemonByNumber := make(map[string]models.Pokemon)
	pokemonByType := make(map[string][]models.Pokemon)

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
		stats := models.Stats{
			Total:   toInt(total),
			Exp:     toInt(total), // Assign Total as Exp in Stats
			HP:      toInt(hp),
			Attack:  toInt(attack),
			Defense: toInt(defense),
			SpAtk:   toInt(spAtk),
			SpDef:   toInt(spDef),
			Speed:   toInt(speed),
		}

		pokemon := models.Pokemon{
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

	limitedPokemons := []models.Pokemon{}
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
	file, _ := os.Create("./pokemon_data/pokedex_data/pokemon_types.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(limitedPokemons)

	fmt.Println("Data saved to pokemon_types.json")
}

func toInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

func arePokemonsEqual(p1 models.Pokemon, p2 models.Pokemon) bool {
	return p1.Name == p2.Name &&
		equalStringSlices(p1.Types, p2.Types) &&
		p1.Number == p2.Number &&
		p1.SubName == p2.SubName &&
		p1.Stats == p2.Stats
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
