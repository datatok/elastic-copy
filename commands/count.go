package commands

import (
	"github.com/ebuildy/elastic-copy/pkg/action"
	"github.com/spf13/cobra"
)

func newCountCmd()  *cobra.Command {
	action := action.NewCountAction(settings)

	cmd := &cobra.Command{
		Use:     "count",
		RunE: func(cmd *cobra.Command, args []string) error {
			action.Run()

			return nil
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringVar(&action.Source, "source", "http://localhost:9200", "source URL")
	flags.StringVar(&action.Query, "query", "", "")
	flags.StringArrayVar(&action.Indices, "indices", nil, "indices to copy")

	return cmd
}
