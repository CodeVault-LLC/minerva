package types

// Default Error Struct for Swagger
type Error struct {
	// Error Code
	Code int `json:"code"`
	// Error Message
	Message string `json:"message"`
}
