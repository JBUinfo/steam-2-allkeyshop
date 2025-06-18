package s2a

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "s2a",
		Short: "Steam wishlist to Allkeyshop",
		Args:  cobra.ArbitraryArgs,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(
		newWishlistCmd(),
	)

	return rootCmd
}
