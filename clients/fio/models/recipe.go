package models

import (
	_ "embed"
	"encoding/json"
)

//go:embed static/recipes.json
var recipeBytes []byte

type Recipe struct {
	BuildingTicker     string        `json:"BuildingTicker"`
	RecipeName         string        `json:"RecipeName"`
	StandardRecipeName string        `json:"StandardRecipeName"`
	Inputs             []InputOutput `json:"Inputs"`
	Outputs            []InputOutput `json:"Outputs"`
	TimeMs             int           `json:"TimeMs"`
}

type InputOutput struct {
	Ticker string `json:"Ticker"`
	Amount int    `json:"Amount"`
}

func StaticRecipeList() []Recipe {
	out := []Recipe{}
	json.Unmarshal(recipeBytes, &out)
	return out
}
