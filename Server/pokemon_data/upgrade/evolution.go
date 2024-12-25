package main

import (
	"fmt"
	"strconv"

	models "Net-Centric--Pokemon/Server/pokemon_data/models"
)

// FindEvolutionChains identifies evolution chains for Pokémon.
func FindEvolutionChains(pokemons []models.Pokemon, pokemonByNumber map[string]models.Pokemon, rangeVal int) map[string][]models.Pokemon {
	evolutionChains := make(map[string][]models.Pokemon)

	for _, pokemon := range pokemons {
		chain := []models.Pokemon{pokemon}
		num, err := strconv.Atoi(pokemon.Number)
		if err != nil {
			fmt.Printf("Error parsing Pokémon number %s: %v\n", pokemon.Number, err)
			continue
		}

		for i := 1; i <= rangeVal; i++ {
			if len(chain) >= 3 {
				break
			}
			nextNumStr := fmt.Sprintf("%04d", num+i)
			nextPokemon, ok := pokemonByNumber[nextNumStr]

			// Add only Pokémon with matching types
			if ok && AreSameTypes(pokemon.Types, nextPokemon.Types) {
				chain = append(chain, nextPokemon)
			}
		}

		if len(chain) > 1 {
			evolutionChains[pokemon.Name] = chain
		}
	}
	return evolutionChains
}

// AreSameTypes checks if two Pokémon share compatible types.
func AreSameTypes(types1, types2 []string) bool {
	typeSet := make(map[string]struct{})
	for _, typ := range types1 {
		typeSet[typ] = struct{}{}
	}
	for _, typ := range types2 {
		if _, found := typeSet[typ]; found {
			return true
		}
	}
	return false
}

// Evolve handles Pokémon evolution based on evolution chains.
func Evolve(p *models.Pokemon, evolutionChains map[string][]models.Pokemon) error {
	chain, ok := evolutionChains[p.Name]
	if !ok {
		return fmt.Errorf("cannot find evolution chain of %s", p.Name)
	}

	// Check evolution possibilities
	if len(chain) > 1 {
		for index, poke := range chain {
			if poke.Number == p.Number {
				if index < len(chain)-1 {
					evolvedPoke := chain[index+1]
					p.Name = evolvedPoke.Name
					p.Types = evolvedPoke.Types
					p.Number = evolvedPoke.Number
					p.Stats = evolvedPoke.Stats
					p.SubName = evolvedPoke.SubName
					p.Level = 1
					p.Stats.Exp = 0
					return nil
				}
				return fmt.Errorf("already reached maximum evolution of %s", p.Name)
			}
		}
	}
	return fmt.Errorf("cannot evolve %s: number not found in evolution chain", p.Name)
}
