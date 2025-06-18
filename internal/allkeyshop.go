package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"io"
	"net/http"
	"os"
	"regexp"
	"s2a/models"
	"strconv"
	"time"
)

const (
	aksBaseURL = "https://www.allkeyshop.com/blog"
	AKSFile    = "aks.json"
)

var regexpNickAndWishlistID = regexp.MustCompile(`href="https://www\.allkeyshop\.com/blog/list/([^/]+)/(\d+)"`)
var regexpProductID = regexp.MustCompile(`(?m)productId\s*=\s*(\d+)`)

type AllKeyShop interface {
	CookieHasValue() bool
	SetCookie(cookie string) error
	ImportToWishlist(fromFile bool, gameTitles []string) error
}

type allKeyShop struct {
	User     models.AKSUser
	Wishlist []models.AKSGame
	NotYet   []models.AKSGame
}

func NewAllKeyShop(sessionCookie string) (AllKeyShop, error) {
	var err error
	aks := &allKeyShop{}
	if _, err = os.Stat(AKSFile); err != nil {
		if err = os.WriteFile(AKSFile, []byte("{}"), 0644); err != nil {
			return aks, err
		}
	}

	currentData, err := aks.readJSON()
	if err != nil {
		return aks, err
	}

	aks.User = currentData.User
	aks.Wishlist = currentData.Wishlist
	aks.NotYet = currentData.NotYet

	if sessionCookie != "" {
		err = aks.SetCookie(sessionCookie)
		if err != nil {
			return aks, err
		}
	}

	return aks, nil
}

func (a *allKeyShop) CookieHasValue() bool {
	return a.User.Cookie != ""
}

func (a *allKeyShop) SetCookie(cookie string) error {
	a.User.Cookie = cookie
	return a.storeJSON()
}

func (a *allKeyShop) ImportToWishlist(fromFile bool, gameTitles []string) error {
	alreadyOnWishlist := make(map[string]bool)
	for _, game := range a.Wishlist {
		alreadyOnWishlist[game.Name] = true
	}

	alreadyOnNotYet := make(map[string]bool)
	for _, game := range a.NotYet {
		alreadyOnNotYet[game.Name] = true
	}

	if fromFile {
		for _, gameTitle := range a.NotYet {
			gameTitles = append(gameTitles, gameTitle.Name)
		}
	}

	for _, gameTitle := range gameTitles {
		slugName := slug.Make(gameTitle)
		if alreadyOnWishlist[slugName] {
			continue
		}

		url := fmt.Sprintf("%s/buy-%s-cd-key-compare-prices/", aksBaseURL, slugName)
		time.Sleep(2 * time.Second)
		body, _, err := a.makeRequest(url, nil)
		if err != nil {
			if !alreadyOnNotYet[slugName] {
				a.NotYet = append(a.NotYet, models.AKSGame{Name: slugName})
			}
			fmt.Printf("WARNING - %s - %s \nError on request - %s\n", gameTitle, slugName, err.Error())
			continue
		}

		matches := regexpProductID.FindStringSubmatch(string(body))
		if len(matches) < 2 {
			if !alreadyOnNotYet[slugName] {
				a.NotYet = append(a.NotYet, models.AKSGame{Name: slugName})
			}
			fmt.Printf("WARNING - %s - %s - couldn't find productId\n", gameTitle, slugName)
			continue
		}

		productId, err := strconv.Atoi(matches[1])
		if err != nil {
			if !alreadyOnNotYet[slugName] {
				a.NotYet = append(a.NotYet, models.AKSGame{Name: slugName})
			}
			fmt.Printf("WARNING - %s - couldn't convert productId to int - %s\n", slugName, err.Error())
			continue
		}

		err = a.addGameToWishlist(models.AKSGame{
			Name: gameTitle,
			ID:   productId,
		})

		if err != nil {
			if !alreadyOnNotYet[slugName] {
				a.NotYet = append(a.NotYet, models.AKSGame{Name: slugName, ID: productId})
			}
			fmt.Printf("WARNING - %s - couldn't add to wishlist - %s\n", slugName, err.Error())
			continue
		}

		fmt.Printf("SAVED! - %s\n", slugName)

		a.Wishlist = append(a.Wishlist, models.AKSGame{Name: slugName, ID: productId})
		if alreadyOnNotYet[slugName] {
			a.removeFromNotYet(slugName)
		}

		err = a.storeJSON()
		if err != nil {
			return err
		}
	}

	return a.storeJSON()
}

func (a *allKeyShop) addGameToWishlist(game models.AKSGame) error {
	if a.User.Nickname == "" || a.User.WishlistID == 0 {
		err := a.updateNicknameAndWishlistID()
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%s/wp-admin/admin-ajax.php?action=akswl_add_game&id=%d&normalisedName=%d", aksBaseURL, a.User.WishlistID, game.ID)
	_, statusCode, err := a.makeRequest(url, nil)

	if statusCode == http.StatusTeapot {
		return nil
	}

	return err
}

func (a *allKeyShop) updateNicknameAndWishlistID() error {
	url := fmt.Sprintf("%s/profile/wishlist/", aksBaseURL)
	body, _, err := a.makeRequest(url, nil)
	if err != nil {
		return err
	}

	matches := regexpNickAndWishlistID.FindStringSubmatch(string(body))
	if len(matches) < 3 {
		return errors.New("couldn't find your wishlist")
	}

	nickname := matches[1]
	wishlistID, err := strconv.Atoi(matches[2])
	if err != nil {
		return fmt.Errorf("error converting wishlist to int: %s", err.Error())
	}

	a.User.Nickname = nickname
	a.User.WishlistID = wishlistID
	return a.storeJSON()
}

func (a *allKeyShop) removeFromNotYet(nameToRemove string) {
	for i, item := range a.NotYet {
		if item.Name == nameToRemove {
			a.NotYet = append(a.NotYet[:i], a.NotYet[i+1:]...)
		}
	}
}

func (a *allKeyShop) makeRequest(url string, bodyReq io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, bodyReq)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %s", err.Error())
	}

	req.AddCookie(&http.Cookie{
		Name:  "wordpress_logged_in_9aae1317051b689fdd8093cf69c60dae",
		Value: a.User.Cookie,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf("URL: %s - status code error: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error reading response: %s", err.Error())
	}

	return body, resp.StatusCode, nil
}

func (a *allKeyShop) storeJSON() error {
	data, err := json.MarshalIndent(models.AKSData{
		User:     a.User,
		Wishlist: a.Wishlist,
		NotYet:   a.NotYet,
	}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(AKSFile, data, 0644)
}

func (a *allKeyShop) readJSON() (models.AKSData, error) {
	var aks models.AKSData
	data, err := os.ReadFile(AKSFile)
	if err != nil {
		return aks, err
	}

	err = json.Unmarshal(data, &aks)
	if err != nil {
		return aks, fmt.Errorf("error decoding %s: %s", AKSFile, err.Error())
	}

	return aks, nil
}
