package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/h2oai/steamY/lib/yarn"
	"github.com/h2oai/steamY/master"
	"github.com/spf13/cobra"
)

const (
	steam = "steam"
)

func Run(version, buildDate string) {
	cmd := Steam(version, buildDate, os.Stdout, os.Stdin, ioutil.Discard)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func Steam(version, buildDate string, stdout, stderr, trace io.Writer) *cobra.Command {
	c := &context{
		version:   version,
		buildDate: buildDate,
	}

	cmd := &cobra.Command{
		Use:               steam,
		Short:             fmt.Sprintf("%s v%s build %s: Command Line Interface to Steam", steam, version, buildDate),
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(
		start(c), // temporary; will not be accessible from the CLI in the future
		stop(c),  // temporary; will not be accessible from the CLI in the future
		serve(c),
	)
	return cmd
}

type context struct {
	version   string
	buildDate string
}

func newCmd(c *context, help string, run func(c *context, args []string)) *cobra.Command {
	doc, err := parseHelp(help)
	if err != nil {
		log.Fatalln("Could not parse help:", err)
	}
	cmd := &cobra.Command{
		Use:   doc.Usage,
		Short: doc.Short,
		Long:  doc.Long,
	}
	if run != nil {
		cmd.Run = func(cmd *cobra.Command, args []string) {
			run(c, args)
		}
	}
	return cmd
}

type Doc struct {
	Usage string
	Short string
	Long  string
}

func parseHelp(text string) (*Doc, error) {
	d := strings.SplitN(strings.TrimSpace(text), "\n", 3)
	if len(d) != 3 {
		return nil, fmt.Errorf("Expected usage, short, long; found %d tokens: %s", len(d), text)
	}
	return &Doc{
		d[0],
		d[1],
		d[1] + "\n\n" + d[2],
	}, nil
}

//
// Commands
//

var serveHelp = `
serve [agent-type]
Lauch a new service.
Examples:

    $ steam serve master
`

func serve(c *context) *cobra.Command {
	cmd := newCmd(c, serveHelp, nil)
	cmd.AddCommand(serveMaster(c))
	return cmd
}

var serveMasterHelp = `
master
Launch the Steam master.
Examples:

    $ steam serve master
`

func serveMaster(c *context) *cobra.Command {
	var webAddress string
	var workingDirectory string
	var enableProfiler bool

	opts := master.DefaultOpts

	cmd := newCmd(c, serveMasterHelp, func(c *context, args []string) {
		master.Run(c.version, c.buildDate, &master.Opts{
			webAddress,
			workingDirectory,
			enableProfiler,
		})
	})

	cmd.Flags().StringVar(&webAddress, "web-address", opts.WebAddress, "Web server address.")
	cmd.Flags().StringVar(&workingDirectory, "working-directory", opts.WorkingDirectory, "Working directory for application files.")
	cmd.Flags().BoolVar(&enableProfiler, "profile", opts.EnableProfiler, "Enable Go profiler")
	return cmd

}

var startHelp = `
start [resource-type]
Start a new resource.
Examples:

    $ steam start cloud
`

func start(c *context) *cobra.Command {
	cmd := newCmd(c, startHelp, nil)
	cmd.AddCommand(startCloud(c))
	return cmd
}

var startCloudHelp = `
cloud
Start a new cloud using the specified H2O package.
Examples:

Start a 4 node H2O 3.2.0.9 cloud

    $ steam start --size=4
`

func startCloud(c *context) *cobra.Command {
	var size int

	cmd := newCmd(c, startCloudHelp, func(c *context, args []string) {
		if len(args) != 1 {
			log.Fatalln("Incorrect number of arguments. See 'steam help start cloud'.")
		}

		name := args[0]
		// version := args[1] // FIXME

		// --- add additional args here ---

		yarn.StartCloud(name, size)

	})
	cmd.Flags().IntVar(&size, "size", 1, "The number of nodes to provision.")
	return cmd
}

var stopHelp = `
stop [resource-type]
Stop the specified resource.
Examples:

    $ steam stop cloud cloud42
`

func stop(c *context) *cobra.Command {
	cmd := newCmd(c, stopHelp, nil)
	cmd.AddCommand(stopCloud(c))
	return cmd
}

var stopCloudHelp = `
cloud [cloud-id]
Stop a cloud.
Examples:

    $ steam stop cloud 1457562501251_0543
`

func stopCloud(c *context) *cobra.Command {
	var force bool

	cmd := newCmd(c, stopCloudHelp, func(c *context, args []string) {
		if len(args) != 1 {
			log.Fatalln("Missing cloud-name. See 'steam help stop cloud'.")
		}

		name := args[0]
		// --- add additional args here ---

		yarn.StopCloud(name)

	})

	cmd.Flags().BoolVar(&force, "force", false, "Force-kill all H2O instances in the cloud")

	return cmd
}
