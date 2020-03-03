package commands

import (
	"github.com/ebuildy/elastic-copy/pkg/action"
	"github.com/spf13/cobra"
)

func newSamplerCmd()  *cobra.Command {
	action := action.NewSampleAction(settings)

	cmd := &cobra.Command{
		Use:     "sample",
		RunE: func(cmd *cobra.Command, args []string) error {
			action.Run()

			return nil
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringVar(&action.Target, "target", "http://localhost:9200", "target URL")
	flags.StringVar(&action.Index, "index", "", "")
	flags.IntVar(&action.Count, "count", 100, "")

	return cmd
}
