package commands

import (
	"github.com/ebuildy/elastic-copy/pkg/action"
	"github.com/spf13/cobra"
)

func newAliasCmd() *cobra.Command {
	action := action.NewAliasAction(settings)

	cmd := &cobra.Command{
		Use: "alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			action.Run()

			return nil
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringVar(&action.Source, "source", "http://localhost:9200", "source URL")
	flags.StringVar(&action.Target, "target", "", "target URL")
	flags.StringVar(&action.IndicesFilter, "filter", "", "regexp to filter indices")

	return cmd
}
