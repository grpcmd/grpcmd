package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/grpcmd/grpcmd/internal/grpcmd"

	"github.com/spf13/cobra"
)

var flagProtoFiles []string
var flagProtoPaths []string

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
		if len(flagProtoFiles) > 0 {
			err := grpcmd.SetFileSource(flagProtoFiles, flagProtoPaths)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while reading proto files:\n\t%s\n", err)
				return grpcmd.ExitError{Code: 1}
			}
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
		} else if len(args) >= 3 {
			headers, data, err := parseRequestFromArgs(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while parsing request from args:\n\t%s\n", err)
				return grpcmd.ExitError{Code: 1}
			}
			err = grpcmd.Call(args[1], data, headers)
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
			if len(flagProtoFiles) > 0 {
				err := grpcmd.SetFileSource(flagProtoFiles, flagProtoPaths)
				if err != nil {
					return emptySlice, directive
				}
			} else {
				err := grpcmd.Connect(args[0])
				if err != nil {
					return emptySlice, directive
				}
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

func parseRequestFromArgs(args []string) ([]string, string, error) {
	var headers []string
	var data strings.Builder
	seenJsonStart := false

	for i := 2; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		if strings.HasPrefix(arg, "{") {
			seenJsonStart = true
		}
		if seenJsonStart {
			data.WriteString(arg)
		} else {
			if strings.IndexByte(arg, ':') == -1 {
				return nil, "", fmt.Errorf("malformed header: missing colon: \"%v\"", arg)
			}
			headers = append(headers, arg)
		}
	}

	return headers, data.String(), nil
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

	rootCmd.Flags().StringSliceVarP(&flagProtoFiles, "protos", "p", nil, "comma-separated list of proto files.")
	rootCmd.Flags().StringSliceVarP(&flagProtoPaths, "paths", "P", nil, "comma-separated list of proto import paths.")
}
