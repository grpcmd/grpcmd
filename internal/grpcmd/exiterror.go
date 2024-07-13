package grpcmd

import (
	"fmt"
)

type ExitError struct {
	Code int
}

func (err ExitError) Error() string {
	return fmt.Sprintf("exit code %v", err.Code)
}

type GrpcStatusExitError struct {
	Code int
}

func (err GrpcStatusExitError) Error() string {
	return fmt.Sprintf("grpc status exit code %v", err.Code)
}
