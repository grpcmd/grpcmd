package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/grpcmd/grpcmd/internal/grpcmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                   "grpc <address> [<method> [<json>...]]",
	Short:                 "A simple, easy-to-use, and developer-friendly CLI tool for gRPC.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		defer grpcmd.Free()
		if len(args) == 0 {
			cmd.Root().Help()
			return
		}
		err := grpcmd.Connect(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while connecting to %s:\n\t%s\n", args[0], err)
			os.Exit(1)
		}
		if len(args) == 1 {
			output, err := grpcmd.ServicesMethodsOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while fetching services and methods from %s:\n\t%s\n", args[0], err)
				os.Exit(1)
			}
			fmt.Println(output)
		} else if len(args) == 2 {
			output, err := grpcmd.DescribeMethod(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while describing method %s:\n\t%s\n", args[1], err)
				os.Exit(1)
			}
			fmt.Println(output)
		} else if len(args) == 3 {
			err := grpcmd.Call(args[1], args[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while calling method %s:\n\t%s\n", args[1], err)
				os.Exit(1)
			}
		} else if len(args) > 3 {
			var data strings.Builder
			for i := 2; i < len(args); i++ {
				data.WriteString(args[i])
			}
			err := grpcmd.Call(args[1], data.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while calling method %s:\n\t%s\n", args[1], err)
				os.Exit(1)
			}
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		emptySlice := []string{}
		directive := cobra.ShellCompDirectiveNoFileComp
		if len(args) == 0 {
			// Do not provide completions for the first argument (hostname:port).
			return emptySlice, directive
		} else if len(args) == 1 {
			defer grpcmd.Free()
			err := grpcmd.Connect(args[0])
			if err != nil {
				return emptySlice, directive
			}
			methods, err := grpcmd.Methods()
			if err != nil {
				return emptySlice, directive
			}
			return methods, directive
		}
		return emptySlice, directive
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(serverCmd)
}
