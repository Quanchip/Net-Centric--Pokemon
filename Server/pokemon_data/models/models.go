// pokemon_data/models/models.go
package models

type Pokemon struct {
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	Number  string   `json:"number"`
	SubName string   `json:"sub_name"`
	Stats   Stats    `json:"stats"`
	EV      float64  `json:"ev"`    // Effort Value (EV) for stats calculation
	Level   int      `json:"level"` // Level of the Pokémon, starting at 0
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
