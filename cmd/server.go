package cmd

import (
	"fmt"
	"os"

	"github.com/grpcmd/grpcmd/internal/grpcmd/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:                   "server [<address>]",
	Short:                 "Starts a gRPC server serving grpcmd.GrpcmdService with reflection",
	DisableFlagsInUseLine: true,
	Args:                  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var address string
		if len(args) > 0 {
			address = args[0]
		} else {
			address = ":50051"
		}
		err := server.Run(address)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while running server:\n\t%s\n", err)
			os.Exit(1)
		}
	},
}
