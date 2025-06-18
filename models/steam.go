package models

type SteamWishlistResponse struct {
	Response struct {
		Items []struct {
			Appid     int `json:"appid"`
			Priority  int `json:"priority"`
			DateAdded int `json:"date_added"`
		} `json:"items"`
	} `json:"response"`
}

type SteamGameDataListResponse struct {
	Response struct {
		StoreItems []struct {
			Id    int    `json:"id"`
			Name  string `json:"name"`
			Appid int    `json:"appid"`
		} `json:"store_items"`
	} `json:"response"`
}

type SteamGame struct {
	Title string `json:"title"`
	AppID int    `json:"appID"`
}

type SteamUser struct {
	SteamKey    string `json:"steamKey"`
	SteamUserID int    `json:"steamUserID"`
}

type SteamData struct {
	User     SteamUser   `json:"user"`
	Wishlist []SteamGame `json:"wishlist"`
}
