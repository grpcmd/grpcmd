package grpcmd

import (
	"fmt"
	"strconv"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcXEventHandler struct {
	grpcurl.DefaultEventHandler
}

func (h *GrpcXEventHandler) OnResolveMethod(md *desc.MethodDescriptor) {}

func (h *GrpcXEventHandler) OnSendHeaders(md metadata.MD) {}

func (h *GrpcXEventHandler) OnReceiveHeaders(md metadata.MD) {
	if md.Len() > 0 {
		fmt.Fprintln(h.Out, grpcurl.MetadataToString(md))
		fmt.Fprintln(h.Out)
	}
}

func (h *GrpcXEventHandler) OnReceiveResponse(resp proto.Message) {
	h.NumResponses++
	if respStr, err := h.Formatter(resp); err != nil {
		fmt.Fprintf(h.Out, "Error while formatting response message #%d:\n\t%s\n", h.NumResponses, err)
	} else {
		fmt.Fprintln(h.Out, respStr)
		fmt.Fprintln(h.Out)
	}
}

func (h *GrpcXEventHandler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	h.Status = stat
	fmt.Fprintln(h.Out, "status-code: "+strconv.Itoa(int(stat.Code()))+" ("+stat.Code().String()+")")
	if md.Len() > 0 {
		fmt.Fprintln(h.Out, grpcurl.MetadataToString(md))
	}
}
