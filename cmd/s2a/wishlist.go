package s2a

import "github.com/spf13/cobra"

func newWishlistCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "wishlist",
		Short: "Interact with your Allkeyshop wishlist",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(
		newWishlistImportCmd(),
	)

	return rootCmd
}
