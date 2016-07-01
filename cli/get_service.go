package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var getServiceHelp = `
service [serviceId]
View detailed service information.
Examples:

	$ steam get service 43
`

func getService(c *context) *cobra.Command {
	cmd := newCmd(c, getServiceHelp, func(c *context, args []string) {
		if len(args) != 1 {
			log.Fatalln("Invalid usage. See 'steam help get service'.")
		}

		// -- Args --

		serviceId, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Invalid usage for serviceId %s: expecting int: %v", args[0], err)
		}

		// -- Execution --

		service, err := c.remote.GetScoringService(serviceId)
		if err != nil {
			log.Fatalln(err)
		}

		model, err := c.remote.GetModel(service.ModelId)
		if err != nil {
			log.Fatalln(err) //FIXME format error
		}

		// -- Formatting --

		base := []string{
			fmt.Sprintf("STATE:\t%s", service.State),
			fmt.Sprintf("PORT:\t%d", service.Port),
			fmt.Sprintf("ID:\t%d", service.Id),
			fmt.Sprintf("AGE:\t%s", fmtAgo(service.CreatedAt)),
		}
		c.printt(fmt.Sprintf("\tScoring Service on Model %d", service.ModelId), base)

		fmt.Println("MODEL")
		ms := []string{
			fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s",
				model.Name,
				model.Id,
				model.Algorithm,
				model.DatasetName,
				model.ResponseColumnName,
				fmtAgo(model.CreatedAt),
			),
		}
		c.printt("NAME\tID\tALGO\tDATASET\tTARGET\tAGE", ms)

	})

	return cmd
}
