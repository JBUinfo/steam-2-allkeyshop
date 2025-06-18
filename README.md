# Steam-2-Allkeyshop
Import your Steam wishlist to Allkeyshop

# Build
1. Install [Golang](https://go.dev/dl/)
2. Open a terminal, go to the project folder and execute:

   `go build -o s2a ./cmd/main.go`
   or
   `go build -o s2a.exe ./cmd/main.go`

# Install
Download your version from [Releases](https://github.com/JBUinfo/YT-chrome-bookmarks-2-MP4/releases)

# Instructions
1. Open a terminal in the project folder en execute:

    `s2a wishlist import -s <ALLKEYSHOP_COOKIE_SESSION> -i <STEAM_USER_ID> -k <STEAM_API_KEY> -r`

After the first run, you will see that there are several games that have failed to be added.
This is because the Steam and Allkeyshop names do not match or, Allkeyshop has not yet indexed these games.
The failed games have been stored in the `aks.json` file in the `notYet` property.
If you have changed any file names or Allkeyshop has updated its list of games, you can re-run the program without reloading the Steam list:

    `s2a wishlist import -f`

# Result
You will have most of the games loaded in your wishlist list in your allkeyshop account.

# Notes
- You can find your <ALLKEYSHOP_COOKIE_SESSION> in the cookie `wordpress_logged_in_9aae1317051b689fdd8093cf69c60dae`.

- You can find your <STEAM_USER_ID> here: [SteamIDFinder](https://www.steamidfinder.com/).

- You can find your <STEAM_API_KEY> here: [SteamApiKey](https://steamcommunity.com/dev/apikey).
