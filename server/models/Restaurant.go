package model

// Restaurant act as model for database
type Restaurant struct {
	RestaurantID int32   `json:"restaurantID,omitempty"`
	Name         string  `json:"name,omitempty"`
	URL          string  `json:"url,omitempty"`
	Cuisines     string  `json:"cuisines,omitempty"`
	Image        string  `json:"image,omitempty"`
	Address      string  `json:"address,omitempty"`
	City         string  `json:"city,omitempty"`
	Rating       float32 `json:"rating,omitempty"`
	IsVeg        bool    `json:"veg,omitempty"`
}
