package commands

import (
	"fmt"
	"github.com/ebuildy/elastic-copy/internal/version"
	"github.com/spf13/cobra"
	"io"
	"text/template"
)

type versionOptions struct {
	short    bool
	template string
}


func newVersionCmd(out io.Writer)  *cobra.Command {
	o := &versionOptions{}

	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"gen"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&o.short, "short", false, "print the version number")
	f.StringVar(&o.template, "template", "", "template for version string format")
	f.BoolP("client", "c", true, "display client version information")
	f.MarkHidden("client")

	return cmd
}

func (o *versionOptions) run(out io.Writer) error {
	if o.template != "" {
		tt, err := template.New("_").Parse(o.template)
		if err != nil {
			return err
		}
		return tt.Execute(out, version.Get())
	}
	fmt.Fprintln(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	v := version.Get()
	if short {
		if len(v.GitCommit) >= 7 {
			return fmt.Sprintf("%s+g%s", v.Version, v.GitCommit[:7])
		}
		return version.GetVersion()
	}
	return fmt.Sprintf("%#v", v)
}