# grpcmd - The "grpc" command.

[![Go Report Card](https://goreportcard.com/badge/github.com/grpcmd/grpcmd)](https://goreportcard.com/report/github.com/grpcmd/grpcmd)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/grpcmd/grpcmd)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/grpcmd/grpcmd/test.yml)
![GitHub Release](https://img.shields.io/github/v/release/grpcmd/grpcmd)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads-pre/grpcmd/grpcmd/total)
![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads-pre/grpcmd/grpcmd/latest/total)

grpcmd is a simple, easy-to-use, and developer-friendly CLI tool for gRPC.

[*(If you're looking for a gRPC GUI desktop app, check out: **grpcmd-gui**)*](https://github.com/grpcmd/grpcmd-gui)

## Demo
![Demo](./demo.svg)

## Feature Overview
 - View Available Services and Methods (using Reflection)
    - `grpc :50051`
 - Describe a gRPC Method
    - `grpc :50051 UnaryMethod`
        - Supports short method names with error on ambiguity.
        - Supports shell completion for the method argument.
        - Outputs a description of the method and its request/response messages.
 - Call a gRPC Method
    - `grpc :50051 UnaryMethod '{"name":"Bob"}'`
        - Supports streaming both requests and responses.
        - Outputs headers, data, and trailers.
        - On gRPC error, `the exit code = 64 + the gRPC status code`.
 - Start a gRPC Server (Implements `grpcmd.GrpcmdService` with Reflection)
    - `grpc server :50051`


## Installation

### Homebrew
    brew install grpcmd/tap/grpcmd
The above command will install the `grpcmd` package from the tap [`grpcmd/tap`](https://github.com/grpcmd/homebrew-tap).

### Binary
You can also download the binary files for macOS, Linux, and Windows from the [Releases](https://github.com/grpcmd/grpcmd/releases) page.

## Usage

### Start a gRPC Server

    grpc server

Output:

    Listening on address: :50051

    Try running:
            grpc :50051 UnaryMethod '{"name": "Bob"}'

This will start a gRPC server (with Reflection) implementing the [`GrpcmdService`](https://github.com/grpcmd/grpcmd/blob/main/proto/grpcmd_service.proto). The methods within this service will also attach the incoming headers as outgoing headers.

#### Options

##### Address
By default, the gRPC server listens for requests on the address `:50051`. To listen on a different address, simply pass it in as the first argument. For example:

    grpc server localhost:12345

### List Available Services and Methods
    grpc :50051

Output:

    grpcmd.GrpcmdService
        BidirectionalStreamingMethod
        ClientStreamingMethod
        ServerStreamingMethod
        UnaryMethod

    grpc.reflection.v1.ServerReflection
        ServerReflectionInfo

    grpc.reflection.v1alpha.ServerReflection
        ServerReflectionInfo

### Describe a Method
    grpc :50051 UnaryMethod

Output:

    rpc UnaryMethod ( .GrpcmdRequest ) returns ( .GrpcmdResponse );

    message GrpcmdRequest {
      string name = 1;
    }

    message GrpcmdResponse {
      string message = 1;
    }

    GrpcmdRequest Template:
    {
      "name": ""
    }

### Call a Unary Method
    grpc :50051 UnaryMethod '{"name": "Bob"}'

Output:

    content-type: application/grpc

    {
      "message": "Hello, Bob!"
    }

    status-code: 0 (OK)

### Call a Client Streaming Method
    grpc :50051 ClientStreamingMethod '{"name": "Bob"}{"name": "Alice"}'

or

    grpc :50051 ClientStreamingMethod '{"name": "Bob"}' '{"name": "Alice"}'

Output:

    content-type: application/grpc

    {
      "message": "Hello, Bob + Alice!"
    }

    status-code: 0 (OK)

### Call a Server Streaming Method
    grpc :50051 ServerStreamingMethod '{"name": "Bob"}'

Output:

    content-type: application/grpc

    {
      "message": "Hello, "
    }

    {
      "message": "B"
    }

    {
      "message": "o"
    }

    {
      "message": "b"
    }

    {
      "message": "!"
    }

    status-code: 0 (OK)

### Call a Bidirectional Streaming Method
    grpc :50051 BidirectionalStreamingMethod '{"name": "Bob"}{"name": "Alice"}'

or

    grpc :50051 BidirectionalStreamingMethod '{"name": "Bob"}' '{"name": "Alice"}'

Output:

    content-type: application/grpc

    {
      "message": "Hello, Bob!"
    }

    {
      "message": "Hello, Alice!"
    }

    status-code: 0 (OK)

## Additional Documentation

### Proto Files
If you want to use `.proto` files instead of Reflection, you can pass one or more comma-separated file locations to the `--protos` (shorthand: `-p`) flag. For example:

    grpc --protos a.proto --protos b.proto,c.proto :50051

### Proto Import Paths
If your `.proto` files contain import statements, you'll likely want to set the search paths for the imports to work properly. To do this, you can pass one or more comma-separated directory locations to the `--paths` (shorthand: `-P`) flag. For example:

    grpc --protos a.proto --paths ../protos/ :50051

### Sending an Empty Request
To send an empty request, simply pass an empty JSON argument for the request data. For example:

    grpc :50051 UnaryMethod {}

### Sending Additional Request Headers (Metadata)
To send a request with additional headers, simply pass one or more `key: value` arguments before the request data. For example:

    grpc :50051 UnaryMethod 'custom-header: custom-value' {}

or

    grpc :50051 UnaryMethod header-1:value-1 header-2:value-2 {}

### Exit Codes
On success, the exit code will be `0`. However, if there is a non-OK gRPC status code in the response, the exit code will be equal to `64 + the gRPC status code`. Other application errors will have exit codes less than `64`. For example, failure to connect to the provided address will result in an exit code of `1`.

### Disambiguating Short Method Names
grpcmd supports taking in short methods names instead of requiring fully-qualified and namespaced method names. In the case where a short method name exists in more that one namespace or package, the tool will throw an error. For example:

    grpc :50051 ServerReflectionInfo

Output:

    Error while describing method ServerReflectionInfo:
        Ambiguous method ServerReflectionInfo. Matching methods:
                grpc.reflection.v1.ServerReflection.ServerReflectionInfo
                grpc.reflection.v1alpha.ServerReflection.ServerReflectionInfo

### Stop Firewall Popup
You may receive a firewall popup when running `grpc server :50051`. For example, this is the case on macOS. The reason for this is that the address `:50051` will listen to port `50051` on **all network interfaces**. To solve this, you can specify the loopback interface in the address. For example:
 - `grpc server localhost:50051`
 - `grpc server 127.0.0.1:50051`

When connecting to the server, you can still start with `grpc :50051` if you prefer.

### Setup Shell Completion

#### Homebrew
If you installed the `grpcmd` package using Homebrew, the shell completion scripts should be installed to their respective directories. If you haven't setup brew completions, follow this [guide](https://docs.brew.sh/Shell-Completion). The guide includes directions for bash, zsh, and fish.

#### Manual
If you want to manually enable shell completion, run the following commands based on your shell. Note: Running the following commands will only affect the current session.

##### Bash

    source <(grpc completion bash)

##### Zsh

    source <(grpc completion zsh)

##### fish

    grpc completion fish | source

##### PowerShell

    bash completion powershell | Out-String | Invoke-Expression
