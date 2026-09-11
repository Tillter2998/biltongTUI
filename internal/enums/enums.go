package enums

import (
	"maps"
	"slices"

	"github.com/Tillter2998/biltongTUI/internal/quicksort"
)

type Ingredients int

const (
	RedWineVinegar Ingredients = iota
	WorcestershireSauce
	Salt
	PepperCorn
	CorianderSeed
	ChiliFlakes
)

var ingredientName = map[Ingredients]string{
	RedWineVinegar:      "Red Wine Vinegar",
	WorcestershireSauce: "Worcestershire Sauce",
	Salt:                "Salt",
	PepperCorn:          "Pepper Corn",
	CorianderSeed:       "Coriander Seed",
	ChiliFlakes:         "Chili Flakes",
}

func (i Ingredients) String() string {
	return ingredientName[i]
}

func AllIngredients() []Ingredients {
	ingredients := slices.Collect(maps.Keys(ingredientName))
	quicksort.Quicksort(ingredients, 0, len(ingredients)-1)

	return ingredients
}
