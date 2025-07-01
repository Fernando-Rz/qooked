package models

// user definition
type User struct {
	UserId      string `json:"id"`
	GroupId     string `json:"groupId"`
	ProfileName string `json:"profileName" binding:"required" `
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
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
	IngredientName string `json:"ingredientName"`
	Amount         string `json:"amount"`
}
