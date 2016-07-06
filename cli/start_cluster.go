package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var startClusterHelp = `
cluster [clusterName] [engineId]
Start a new cluster using the specified H2O package.
Examples:

Start a 4-node H2O 3.2.0.9 cluster

    $ steam start cluster 42 2 --size=4
`

func startCluster(c *context) *cobra.Command {
	var (
		size   int
		memory string
	)

	cmd := newCmd(c, startClusterHelp, func(c *context, args []string) {
		if len(args) != 2 {
			log.Fatalln("Incorrect number of arguments. See 'steam help start cluster'.")
		}

		clusterName := args[0]
		engineId, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("Ivalid usage of engineId %q: expecting integer: %v", args[1], err)
		}

		// --- add additional args here ---

		log.Println("Attempting to start cluster...")
		clusterId, err := c.remote.StartYarnCluster(clusterName, engineId, size, memory, "")
		if err != nil {
			log.Fatalln(err)
		}

		cluster, err := c.remote.GetCluster(clusterId)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Cluster %d started at: %s", clusterId, cluster.Address)

		// if details {
		// 	info, err := c.remote.GetClusterStatus(cluster.Name)
		// 	if err != nil {
		// 		log.Fatalln(err)
		// 	}

		// 	fmt.Printf(
		// 		`
		// Engine: %v
		// Version: %v
		// Total Nodes: %v
		// Total Memory: %v
		// `, info.EngineName, info.EngineVersion, info.Size, info.Memory)
		// }
		// // TODO: name corresponds to id for purpose of stopCluster
	})

	cmd.Flags().IntVarP(&size, "size", "n", 1, "The number of nodes to provision.")
	cmd.Flags().StringVarP(&memory, "memory", "m", "10g", "The max amount of memory to use per node.")

	return cmd
}
