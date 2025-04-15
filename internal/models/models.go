package models

// user definition
type User struct {
	UserId      string `json:"id"`
	GroupId     string `json:"groupId"`
	ProfileName string `json:"profileName"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// recipe definition
type Recipe struct {
	RecipeId     string       `json:"id"`
	UserId       string       `json:"userId"`
	RecipeName   string       `json:"recipeName"`
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
