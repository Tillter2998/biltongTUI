package enums

import "errors"

type Ingredients int

const (
	RedWineVinegar Ingredients = iota
	WorcestershireSauce
	Salt
	PepperCorn
	CorianderSeed
	ChiliFlakes
)

var IngredientName = map[Ingredients]string{
	RedWineVinegar:      "Red Wine Vinegar",
	WorcestershireSauce: "Worcestershire Sauce",
	Salt:                "Salt",
	PepperCorn:          "Pepper Corn",
	CorianderSeed:       "Coriander Seed",
	ChiliFlakes:         "Chili Flakes",
}

func (i Ingredients) String() string {
	return IngredientName[i]
}

func StringToIngredients(value string) (Ingredients, error) {
	for ingredient, ingredientValue := range IngredientName {
		if ingredientValue == value {
			return ingredient, nil
		}
	}
	return 0, errors.New("Could not find matching Ingredient")
}
