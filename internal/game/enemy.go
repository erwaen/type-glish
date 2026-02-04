package game

import (
	"math/rand"
)

type Enemy struct {
	Name        string
	HP          int
	MaxHP       int
	Location    string
	Description string
}

// Predefined enemies
var Enemies = []Enemy{
	{
		Name:        "Goblin",
		HP:          30,
		MaxHP:       30,
		Location:    "The Murky Swamp",
		Description: "A sneaky goblin with a rusty dagger, muttering broken sentences.",
	},
	{
		Name:        "Troll",
		HP:          80,
		MaxHP:       80,
		Location:    "The Whispering Woods",
		Description: "A massive troll with a wooden club. He mocks your grammar mistakes.",
	},
	{
		Name:        "Skeleton",
		HP:          40,
		MaxHP:       40,
		Location:    "The Crypt of Conjugations",
		Description: "A rattling skeleton that speaks only in past tense.",
	},
	{
		Name:        "Dark Wizard",
		HP:          60,
		MaxHP:       60,
		Location:    "The Tower of Tenses",
		Description: "A hooded figure casting spells with perfectly structured incantations.",
	},
	{
		Name:        "Grammar Golem",
		HP:          100,
		MaxHP:       100,
		Location:    "The Lexicon Library",
		Description: "A towering construct made of ancient dictionaries and thesauri.",
	},
	{
		Name:        "Syntax Spider",
		HP:          35,
		MaxHP:       35,
		Location:    "The Web of Words",
		Description: "A giant spider that weaves webs of confusing clauses.",
	},
}

// RandomEnemy returns a copy of a random enemy from the list
func RandomEnemy() *Enemy {
	idx := rand.Intn(len(Enemies))
	enemy := Enemies[idx]
	return &Enemy{
		Name:        enemy.Name,
		HP:          enemy.MaxHP,
		MaxHP:       enemy.MaxHP,
		Location:    enemy.Location,
		Description: enemy.Description,
	}
}

