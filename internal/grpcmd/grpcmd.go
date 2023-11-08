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
	_refClient *grpcreflect.Client
	_refSource grpcurl.DescriptorSource

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

	_refClient = grpcreflect.NewClientAuto(_ctx, _cc)
	deferCall(_refClient.Reset)
	_refSource = grpcurl.DescriptorSourceFromServer(_ctx, _refClient)
	return nil
}

func Services() ([]string, error) {
	if _services != nil {
		return _services, nil
	}
	services, err := grpcurl.ListServices(_refSource)
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
		methods, err := grpcurl.ListMethods(_refSource, s)
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

func ServicesMethodsOutput() (string, error) {
	_, err := Methods()
	if err != nil {
		return "", err
	}
	return _servicesMethodsOutput.String(), nil
}

func DescribeMethod(method string) (string, error) {
	var output strings.Builder
	dsc, err := _refSource.FindSymbol(method)
	if err != nil {
		return "", err
	}
	txt, err := grpcurl.GetDescriptorText(dsc, _refSource)
	if err != nil {
		return "", err
	}
	output.WriteString(txt)
	output.WriteRune('\n')
	output.WriteRune('\n')

	// TODO: it is possible to convert this into an if statement and ok conversion check.
	switch d := dsc.(type) {
	case *desc.MethodDescriptor:
		txt, err = grpcurl.GetDescriptorText(d.GetInputType(), _refSource)
		if err != nil {
			return "", err
		}
		output.WriteString(txt)
		output.WriteRune('\n')
		output.WriteRune('\n')
		txt, err = grpcurl.GetDescriptorText(d.GetOutputType(), _refSource)
		if err != nil {
			return "", err
		}
		output.WriteString(txt)
		output.WriteRune('\n')
		output.WriteRune('\n')

		tmpl := grpcurl.MakeTemplate(d.GetInputType())
		options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
		_, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, _refSource, nil, options)
		if err != nil {
			return "", err
		}
		str, err := formatter(tmpl)
		if err != nil {
			return "", err
		}
		output.WriteString(d.GetInputType().GetName() + " Template:\n")
		output.WriteString(str)
	default:
		return "", errors.New("Descriptor for " + dsc.GetFullyQualifiedName() + " is not a MethodDescriptor.")
	}
	return output.String(), nil
}

func Call(address, data string) error {
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		AllowUnknownFields:    false,
		IncludeTextSeparator:  false,
	}
	rp, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, _refSource, strings.NewReader(data), options)
	if err != nil {
		return err
	}
	h := &GrpcXEventHandler{
		DefaultEventHandler: grpcurl.DefaultEventHandler{
			Out:            os.Stdout,
			Formatter:      formatter,
			VerbosityLevel: 0,
		},
	}
	err = grpcurl.InvokeRPC(_ctx, _refSource, _cc, address, nil, h, rp.Next)
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
		os.Exit(64 + int(h.Status.Code()))
	}
	return nil
}
