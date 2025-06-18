package models

type AKSGame struct {
	Name string `json:"name"`
	ID   int    `json:"id,omitempty"`
}

type AKSUser struct {
	Nickname   string `json:"nickname"`
	Cookie     string `json:"cookie"`
	WishlistID int    `json:"wishlistID"`
}

type AKSData struct {
	User     AKSUser   `json:"user"`
	Wishlist []AKSGame `json:"wishlist"`
	NotYet   []AKSGame `json:"notYet"`
}
