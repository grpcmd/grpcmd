package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/grpcmd/grpcmd/internal/grpcmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "grpc <address> [<method> [<json>...]]",
	Short:         "A simple, easy-to-use, and developer-friendly CLI tool for gRPC.",
	Args:          cobra.MinimumNArgs(0),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		defer grpcmd.Free()
		if len(args) == 0 {
			cmd.Root().Help()
			return nil
		}
		err := grpcmd.Connect(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while connecting to %s:\n\t%s\n", args[0], err)
			return grpcmd.ExitError{Code: 1}
		}
		if len(args) == 1 {
			output, err := grpcmd.ServicesMethodsOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while fetching services and methods from %s:\n\t%s\n", args[0], err)
				return grpcmd.ExitError{Code: 1}
			}
			fmt.Println(output)
		} else if len(args) == 2 {
			output, err := grpcmd.DescribeMethod(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while describing method %s:\n\t%s\n", args[1], err)
				return grpcmd.ExitError{Code: 1}
			}
			fmt.Println(output)
		} else if len(args) == 3 {
			err := grpcmd.Call(args[1], args[2])
			if err != nil {
				if e := new(grpcmd.GrpcStatusExitError); errors.As(err, e) {
					return grpcmd.ExitError{Code: e.Code}
				}
				fmt.Fprintf(os.Stderr, "Error while calling method %s:\n\t%s\n", args[1], err)
				return grpcmd.ExitError{Code: 1}
			}
		} else if len(args) > 3 {
			var data strings.Builder
			for i := 2; i < len(args); i++ {
				data.WriteString(args[i])
			}
			err := grpcmd.Call(args[1], data.String())
			if err != nil {
				if e := new(grpcmd.GrpcStatusExitError); errors.As(err, e) {
					return grpcmd.ExitError{Code: e.Code}
				}
				fmt.Fprintf(os.Stderr, "Error while calling method %s:\n\t%s\n", args[1], err)
				return grpcmd.ExitError{Code: 1}
			}
		}
		return nil
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
			methods, err := grpcmd.NonambiguousMethods()
			if err != nil {
				return emptySlice, directive
			}
			return methods, directive
		}
		return emptySlice, directive
	},
}

func SetBuildInfo(version string, date string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("grpc version " + version + " (" + date[0:10] + ")\n" +
		"https://github.com/grpcmd/grpcmd/releases/v" + version + "\n")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		if e := new(grpcmd.ExitError); errors.As(err, e) {
			os.Exit(e.Code)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(serverCmd)
}
