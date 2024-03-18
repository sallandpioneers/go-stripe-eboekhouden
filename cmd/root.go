package cmd

import "github.com/spf13/cobra"

//nolint:gochecknoglobals // this is the recommended way to use cobra
var RootCmd = &cobra.Command{
	Use:   "go-stripe-eboekhouden",
	Short: "Stripe eboekhouden",
}

func Execute() error {
	return RootCmd.Execute()
}

//nolint:gochecknoinits // this is the recommended way to use cobra
func init() {
	RootCmd.AddCommand(serveCommand)
	RootCmd.AddCommand(migrateCommand)
}
