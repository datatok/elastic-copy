package commands

import (
	"fmt"
	"io"

	"github.com/ebuildy/elastic-copy/pkg/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	settings = cli.New()
)

func NewRootCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "elasticcopy",
		Short:        "Copy data from & to elasticsearch",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			l, err := logrus.ParseLevel(settings.LogLevel)

			if err != nil {
				fmt.Printf("log level \"%s\" is invalid!", settings.LogLevel)
				l = logrus.InfoLevel
			}

			logrus.SetLevel(l)
		},
	}

	cmd.AddCommand(
		newRunCmd(),
		newSamplerCmd(),
		newCountCmd(),
		newAliasCmd(),
		newVersionCmd(out),
	)

	flags := cmd.PersistentFlags()

	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.Parse(args)

	return cmd
}
