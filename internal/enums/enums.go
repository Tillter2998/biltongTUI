package enums

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

// var IngredientName = map[Ingredients]values{
// 	redWineVinegar:      {value: "Red Wine Vinegar", checked: true},
// 	worcestershireSauce: {value: "Worcestershire Sauce", checked: true},
// 	salt:                {value: "Salt", checked: true},
// 	pepperCorn:          {value: "Pepper Corn", checked: true},
// 	corianderSeed:       {value: "Coriander Seed", checked: true},
// 	chiliFlakes:         {value: "Chili Flakes"},
// }
