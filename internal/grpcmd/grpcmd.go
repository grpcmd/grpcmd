package grpcmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ctx       context.Context
	_cc        *grpc.ClientConn
	_dscSource grpcurl.DescriptorSource

	_services              []string
	_methods               []string
	_servicesMethodsOutput strings.Builder

	_freeQueue []func()
)

func deferCall(f func()) {
	_freeQueue = append(_freeQueue, f)
}

func Free() {
	for i := len(_freeQueue) - 1; i >= 0; i-- {
		_freeQueue[i]()
	}
}

func SetFileSource(protoFiles, protoPaths []string) error {
	fileSource, err := grpcurl.DescriptorSourceFromProtoFiles(
		// Deduplication is required because for some reason the following command parses duplicate flags.
		// $ grpc __complete --protos ./proto/grpcmd_service.proto :50051 UnaryMethod
		removeDuplicates(protoPaths),
		removeDuplicates(protoFiles)...,
	)
	if err != nil {
		return err
	}
	_dscSource = fileSource
	return nil
}

func removeDuplicates[T comparable](slice []T) []T {
	unique := make([]T, 0, len(slice))
	seen := make(map[T]bool)

	for _, value := range slice {
		if !seen[value] {
			seen[value] = true
			unique = append(unique, value)
		}
	}

	return unique
}

func Connect(address string) error {
	var cancel context.CancelFunc
	_ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	deferCall(cancel)

	var err error
	_cc, err = grpcurl.BlockingDial(_ctx, "tcp", address, nil)
	if err != nil {
		return err
	}
	deferCall(func() { _cc.Close() })

	if _dscSource == nil {
		refClient := grpcreflect.NewClientAuto(_ctx, _cc)
		deferCall(refClient.Reset)
		refSource := grpcurl.DescriptorSourceFromServer(_ctx, refClient)
		_dscSource = refSource
	}
	return nil
}

func Services() ([]string, error) {
	if _services != nil {
		return _services, nil
	}
	services, err := grpcurl.ListServices(_dscSource)
	if err != nil {
		return nil, err
	}
	_services = services
	return _services, nil
}

func Methods() ([]string, error) {
	if _methods != nil {
		return _methods, nil
	}
	services, err := Services()
	if err != nil {
		return nil, err
	}
	for _, s := range services {
		methods, err := grpcurl.ListMethods(_dscSource, s)
		if err != nil {
			return nil, err
		}
		_methods = append(_methods, methods...)
		_servicesMethodsOutput.WriteString(s)
		_servicesMethodsOutput.WriteRune('\n')
		for _, m := range methods {
			_servicesMethodsOutput.WriteRune('\t')
			_servicesMethodsOutput.WriteString(m[len(s)+1:])
			_servicesMethodsOutput.WriteRune('\n')
		}
		_servicesMethodsOutput.WriteRune('\n')
	}
	return _methods, nil
}

func NonambiguousMethods() ([]string, error) {
	methods, err := Methods()
	if err != nil {
		return nil, err
	}

	nonambiguousMethods := make([]string, 0, len(methods))
	ambiguousMethods := make(map[string]bool)

	for _, fullyQualifiedName := range methods {
		i := strings.LastIndex(fullyQualifiedName, ".")
		var name string
		if i == -1 {
			name = fullyQualifiedName
		} else {
			name = fullyQualifiedName[i+1:]
		}
		nonambiguousMethods = append(nonambiguousMethods, name)
		if _, ok := ambiguousMethods[name]; ok {
			ambiguousMethods[name] = true
		} else {
			ambiguousMethods[name] = false
		}
	}

	for i, fullyQualifiedName := range methods {
		name := nonambiguousMethods[i]
		if ambiguousMethods[name] {
			nonambiguousMethods[i] = fullyQualifiedName
		}
	}

	return nonambiguousMethods, nil
}

func findFullyQualifiedMethod(method string) (string, error) {
	methods, err := Methods()
	if err != nil {
		return "", err
	}
	matches := make([]string, 0, 1)
	exactMatches := make([]string, 0, 1)
	for _, fullyQualifiedName := range methods {
		if i := strings.Index(fullyQualifiedName, method); i > -1 {
			matches = append(matches, fullyQualifiedName)
		}
		i := strings.LastIndex(fullyQualifiedName, ".")
		name := fullyQualifiedName[i+1:]
		if method == name {
			exactMatches = append(exactMatches, fullyQualifiedName)
		}
	}
	if len(matches) == 0 {
		return "", errors.New("No matching method for: " + method)
	} else if len(matches) == 1 {
		return matches[0], nil
	} else if len(exactMatches) == 1 {
		return exactMatches[0], nil
	} else {
		var text strings.Builder
		text.WriteString("Ambiguous method ")
		text.WriteString(method)
		text.WriteString(". Matching methods:\n")
		for _, m := range matches {
			text.WriteString("\t\t")
			text.WriteString(m)
			text.WriteRune('\n')
		}
		return "", errors.New(text.String())
	}
}

func ServicesMethodsOutput() (string, error) {
	_, err := Methods()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(_servicesMethodsOutput.String(), "\n"), nil
}

func DescribeMethod(method string) (string, error) {
	fullyQualifiedMethod, err := findFullyQualifiedMethod(method)
	if err != nil {
		return "", err
	}
	var output strings.Builder
	dsc, err := _dscSource.FindSymbol(fullyQualifiedMethod)
	if err != nil {
		return "", err
	}
	txt, err := grpcurl.GetDescriptorText(dsc, _dscSource)
	if err != nil {
		return "", err
	}
	output.WriteString(txt)
	output.WriteRune('\n')
	output.WriteRune('\n')

	if d, ok := dsc.(*desc.MethodDescriptor); ok {
		txt, err = grpcurl.GetDescriptorText(d.GetInputType(), _dscSource)
		if err != nil {
			return "", err
		}
		output.WriteString(txt)
		output.WriteRune('\n')
		output.WriteRune('\n')
		txt, err = grpcurl.GetDescriptorText(d.GetOutputType(), _dscSource)
		if err != nil {
			return "", err
		}
		output.WriteString(txt)
		output.WriteRune('\n')
		output.WriteRune('\n')

		tmpl := grpcurl.MakeTemplate(d.GetInputType())
		options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
		_, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, _dscSource, nil, options)
		if err != nil {
			return "", err
		}
		str, err := formatter(tmpl)
		if err != nil {
			return "", err
		}
		output.WriteString(d.GetInputType().GetName() + " Template:\n")
		output.WriteString(str)
	} else {
		return "", errors.New("Descriptor for " + dsc.GetFullyQualifiedName() + " is not a MethodDescriptor.")
	}
	return output.String(), nil
}

func Call(method, data string, headers []string) error {
	fullyQualifiedMethod, err := findFullyQualifiedMethod(method)
	if err != nil {
		return err
	}
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		AllowUnknownFields:    false,
		IncludeTextSeparator:  false,
	}
	rp, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, _dscSource, strings.NewReader(data), options)
	if err != nil {
		return err
	}
	h := &GrpcmdEventHandler{
		DefaultEventHandler: grpcurl.DefaultEventHandler{
			Out:            os.Stdout,
			Formatter:      formatter,
			VerbosityLevel: 0,
		},
	}
	err = grpcurl.InvokeRPC(_ctx, _dscSource, _cc, fullyQualifiedMethod, headers, h, rp.Next)
	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			h.Status = errStatus
		} else {
			return err
		}
	}
	if h.Status.Code() != codes.OK {
		formattedStatus, err := formatter(h.Status.Proto())
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, formattedStatus)
		return GrpcStatusExitError{Code: 64 + int(h.Status.Code())}
	}
	return nil
}
