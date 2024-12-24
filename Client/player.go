// PLAYER CODE
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

type PlayerRequest struct {
	Name string `json:"name"`
}

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type PlayerResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	player := connect()
	join()

	//
	ok := battlePromt("Go to battle mode?", true)
	if ok {
		fmt.Println("Go to battle!")
		fmt.Printf("Player ID: %d\n", player.ID)
	} else {
		fmt.Println("Are you scared?")
		fmt.Printf("Player ID: %d\n", player.ID)
	}
}

func connect() PlayerResponse {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]

	// Create request body
	requestBody := PlayerRequest{
		Name: name,
	}
	jsonData, _ := json.Marshal(requestBody)

	// Send POST request
	resp, err := http.Post("http://localhost:8080/connect", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return PlayerResponse{}
	}
	defer resp.Body.Close()

	// Read response
	var player PlayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		fmt.Println("Error reading response:", err)
		return PlayerResponse{}
	}

	return player
}

func join() {
	// Make a POST request to join the game
	response, err := http.Post("http://localhost:8080/join", "application/json", nil)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: received status code", response.StatusCode)
		return
	}

	// Decode the response JSON
	var result struct {
		Pokemons []Pokemon `json:"pokemons"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Display the received Pokémon
	fmt.Println("You received the following Pokémon:")
	for _, p := range result.Pokemons {
		fmt.Printf("- %s (Type: %v, HP: %d, ATK: %d, DEF: %d)\n", p.Name, p.Types, p.Stats.HP, p.Stats.Attack, p.Stats.Defense)
	}
}

func battlePromt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
