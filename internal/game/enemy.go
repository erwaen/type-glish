package game

import (
	"math/rand"
)

type Enemy struct {
	Name        string
	HP          int
	MaxHP       int
	Tier        int // 1=easy, 2=medium, 3=hard, 4=boss
	Location    string
	Description string
}

// Predefined enemies
var Enemies = []Enemy{
	{
		Name:        "Goblin",
		HP:          20,
		MaxHP:       20,
		Tier:        1,
		Location:    "The Murky Swamp",
		Description: "A sneaky goblin with a rusty dagger, muttering broken sentences.",
	},
	{
		Name:        "Syntax Spider",
		HP:          25,
		MaxHP:       25,
		Tier:        1,
		Location:    "The Web of Words",
		Description: "A giant spider that weaves webs of confusing clauses.",
	},
	{
		Name:        "Skeleton",
		HP:          30,
		MaxHP:       30,
		Tier:        2,
		Location:    "The Crypt of Conjugations",
		Description: "A rattling skeleton that speaks only in past tense.",
	},
	{
		Name:        "Dark Wizard",
		HP:          40,
		MaxHP:       40,
		Tier:        2,
		Location:    "The Tower of Tenses",
		Description: "A hooded figure casting spells with perfectly structured incantations.",
	},
	{
		Name:        "Troll",
		HP:          50,
		MaxHP:       50,
		Tier:        3,
		Location:    "The Whispering Woods",
		Description: "A massive troll with a wooden club. He mocks your grammar mistakes.",
	},
	{
		Name:        "Grammar Golem",
		HP:          75,
		MaxHP:       75,
		Tier:        4,
		Location:    "The Lexicon Library",
		Description: "A towering construct made of ancient dictionaries and thesauri.",
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
		Tier:        enemy.Tier,
		Location:    enemy.Location,
		Description: enemy.Description,
	}
}

