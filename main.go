package main

import "github.com/grpcmd/grpcmd/cmd"

var (
	version = "0.0.0"
	date    = "0000-00-00T00:00:00Z"
)

func main() {
	cmd.SetBuildInfo(version, date)
	cmd.Execute()
}
