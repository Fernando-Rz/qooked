package models

// recipe definition
type Recipe struct {
    UserID       string
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

// user accounts
type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"` // gotta secure this
}
