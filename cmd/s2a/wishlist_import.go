package s2a

import (
	"fmt"
	"github.com/spf13/cobra"
	"s2a/internal"
)

func newWishlistImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import your Steam wishlist to Allkeyshop",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fromFile := cmd.Flags().Changed(aksFromFile)
			var gameTitles []string
			if !fromFile {
				steamID, err := cmd.Flags().GetInt(steamIDFlag)
				if err != nil {
					return err
				}

				steamKey, err := cmd.Flags().GetString(steamKeyFlag)
				if err != nil {
					return err
				}

				fmt.Println("Reading Steam games from file...")
				steam, err := internal.NewSteam(steamID, steamKey)
				if err != nil {
					return err
				}

				if !steam.SteamIDHasValue() {
					return fmt.Errorf("steamID is required. please use -i <STEAM_USER_ID>")
				}

				if !steam.SteamIDHasValue() {
					return fmt.Errorf("steamKey is required. please use -k <STEAM_API_KEY>")
				}

				shouldRefresh := cmd.Flags().Changed(steamRefreshWishlist)
				if err != nil {
					return err
				}

				if shouldRefresh {
					fmt.Println("Reading Steam games from cloud...")
				}

				wishlist, err := steam.GetWishlist(shouldRefresh)
				if err != nil {
					return err
				}

				for _, game := range wishlist {
					gameTitles = append(gameTitles, game.Title)
				}
			}

			sessionCookie, err := cmd.Flags().GetString(aksSessionCookieFlag)
			if err != nil {
				return err
			}

			fmt.Println("Reading AKS games from file...")
			aks, err := internal.NewAllKeyShop(sessionCookie)
			if err != nil {
				return err
			}

			if !aks.CookieHasValue() {
				return fmt.Errorf("session cookie is required. please use -s <YOUR_COOKIE>")
			}

			fmt.Println("Adding games to AKS...")
			return aks.ImportToWishlist(fromFile, gameTitles)
		},
	}

	cmd.Flags().BoolP(steamRefreshWishlist, "r", false, fmt.Sprintf("Refresh steam wishlist stored in %s", internal.SteamFile))
	cmd.Flags().BoolP(aksFromFile, "f", false, fmt.Sprintf("(Steam is omitted) Try to add the games that failed previously and are stored in 'notYet' from %s", internal.AKSFile))
	cmd.Flags().StringP(aksSessionCookieFlag, "s", "", "(Optional if previously used) Your session cookie")
	cmd.Flags().IntP(steamIDFlag, "i", 0, "(Optional if previously used) Your Steam user ID")
	cmd.Flags().StringP(steamKeyFlag, "k", "", "(Optional if previously used) Your Steam API Key")

	return cmd
}
