package commands

import (
	"github.com/ebuildy/elastic-copy/pkg/action"
	"github.com/ebuildy/elastic-copy/pkg/utils"
	"github.com/spf13/cobra"
)

func newRunCmd()  *cobra.Command {
	action := action.NewRunAction(settings)

	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"gen"},
		RunE: func(cmd *cobra.Command, args []string) error {
			action.Run()

			return nil
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringVar(&action.Source, "source", "http://localhost:9200", "source URL")
	flags.StringVar(&action.Target, "target", utils.TARGET_STDOUT, "target URL")
	flags.StringVar(&action.TargetIndexType, "target_index_type",  "", "target index type (blank = no type)")
	flags.StringVar(&action.Query, "query", "", "")
	flags.StringArrayVar(&action.Indices, "indices", nil, "indices to copy")
	flags.Uint64Var(&action.Count, "count", 0, "0 => all, X => count")
	flags.IntVar(&action.ReadBatchSize, "read_batch", 100, "how many documents to read in one scroll")
	flags.IntVar(&action.WriteBatchSize, "write_batch", 20, "how many documents to send to writer in one batch")
	flags.IntVar(&action.Threads, "threads", 5, "Number of threads in pool")
	flags.StringVar(&action.ForceType, "type-override", "", "")
	flags.BoolVar(&action.FailFast, "fail-fast", true, "")

	return cmd
}
