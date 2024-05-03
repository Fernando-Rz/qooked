package models

// recipe definition
type Recipe struct {
    Name         string       `json:"name"`
    Description  string       `json:"description"`
    Time         RecipeTime   `json:"time"`
    Servings     int          `json:"servings"`
    Ingredients  []Ingredient `json:"ingredients"`
    Instructions []string     `json:"instructions"`
}