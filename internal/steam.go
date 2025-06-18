package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"s2a/models"
)

const (
	steamURLListGamesData = "https://api.steampowered.com/IStoreBrowseService/GetItems/v1/"
	steamURLGetWishlist   = "https://api.steampowered.com/IWishlistService/GetWishlist/v1/?steamid=%d"
	appID                 = "appid"
	SteamFile             = "steam.json"
)

type Steam interface {
	SteamIDHasValue() bool
	SteamKeyHasValue() bool
	SetSteamID(steamID int) error
	SetSteamKey(steamKey string) error
	GetWishlist(shouldRefresh bool) ([]models.SteamGame, error)
	ListGamesData(gameIDs []int) ([]models.SteamGame, error)
}

type steam struct {
	User     models.SteamUser
	Wishlist []models.SteamGame
}

func NewSteam(steamID int, steamKey string) (Steam, error) {
	var err error
	st := &steam{}
	if _, err = os.Stat(SteamFile); err != nil {
		if err = os.WriteFile(SteamFile, []byte("{}"), 0644); err != nil {
			return st, err
		}
	}

	currentData, err := st.readJSON()
	if err != nil {
		return st, err
	}

	st.User = currentData.User
	st.Wishlist = currentData.Wishlist

	if steamID != 0 {
		err = st.SetSteamID(steamID)
		if err != nil {
			return st, err
		}
	}

	if steamKey != "" {
		err = st.SetSteamKey(steamKey)
		if err != nil {
			return st, err
		}
	}

	return st, nil
}

func (s *steam) GetWishlist(shouldRefresh bool) ([]models.SteamGame, error) {
	if !shouldRefresh {
		return s.Wishlist, nil
	}

	urlReq := fmt.Sprintf(steamURLGetWishlist, s.User.SteamUserID)
	resBody, statusCode, err := s.makeRequest(urlReq, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("steam returned status code %d", statusCode)
	}

	var res models.SteamWishlistResponse
	if err = json.Unmarshal(resBody, &res); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %s", err.Error())
	}

	var ids []int
	for _, game := range res.Response.Items {
		ids = append(ids, game.Appid)
	}

	return s.ListGamesData(ids)
}

func (s *steam) SteamIDHasValue() bool {
	return s.User.SteamUserID != 0
}

func (s *steam) SetSteamID(steamID int) error {
	s.User.SteamUserID = steamID
	return s.storeJSON()
}

func (s *steam) SteamKeyHasValue() bool {
	return s.User.SteamKey != ""
}

func (s *steam) SetSteamKey(steamKey string) error {
	s.User.SteamKey = steamKey
	return s.storeJSON()
}

func (s *steam) ListGamesData(gameIDs []int) ([]models.SteamGame, error) {
	if len(gameIDs) == 0 {
		return nil, errors.New("no games found")
	}

	var ids []map[string]int
	for _, id := range gameIDs {
		ids = append(ids, map[string]int{appID: id})
	}

	inputData := map[string]interface{}{
		"ids": ids,
		"context": map[string]interface{}{
			"language":     "english",
			"country_code": "US",
			"steam_realm":  1,
		},
	}

	inputJSONBytes, err := json.Marshal(inputData)
	if err != nil {
		log.Fatalf("Error al codificar JSON: %v", err)
	}

	params := url.Values{}
	params.Set("key", s.User.SteamKey)
	params.Set("input_json", string(inputJSONBytes))

	urlReq := fmt.Sprintf("%s?%s", steamURLListGamesData, params.Encode())

	resBody, statusCode, err := s.makeRequest(urlReq, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("steam returned status code %d", statusCode)
	}

	var res models.SteamGameDataListResponse
	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %s", err.Error())
	}

	uniqueGames := make(map[int]bool)
	for _, game := range s.Wishlist {
		uniqueGames[game.AppID] = true
	}
	
	for _, game := range res.Response.StoreItems {
		if !uniqueGames[game.Appid] {
			s.Wishlist = append(s.Wishlist, models.SteamGame{
				AppID: game.Appid,
				Title: game.Name,
			})
		}
	}

	return s.Wishlist, s.storeJSON()
}

func (s *steam) makeRequest(url string, bodyReq io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, bodyReq)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %s", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error reading response: %s", err.Error())
	}

	return body, resp.StatusCode, nil
}

func (s *steam) storeJSON() error {
	data, err := json.MarshalIndent(models.SteamData{
		User:     s.User,
		Wishlist: s.Wishlist,
	}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(SteamFile, data, 0644)
}

func (s *steam) readJSON() (models.SteamData, error) {
	var st models.SteamData
	data, err := os.ReadFile(SteamFile)
	if err != nil {
		return st, err
	}

	err = json.Unmarshal(data, &st)
	if err != nil {
		return st, fmt.Errorf("error decoding %s: %s", SteamFile, err.Error())
	}

	return st, nil
}
