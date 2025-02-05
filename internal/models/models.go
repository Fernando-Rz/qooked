package models

// recipe definition
type Recipe struct {
	RecipeId     string       `json:"id"`
	UserId       string       `json:"userId"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Time         RecipeTime   `json:"time"`
	Servings     int          `json:"servings"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
}

// time it takes to complete recipe
type RecipeTime struct {
	Prep  string `json:"prep"`
	Cook  string `json:"cook"`
	Total string `json:"total"`
}

// ingredient definition
type Ingredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

// user definition
type User struct {
	Id           string `json:"id"`
	PartitionKey string `json:"partitionKey"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}
