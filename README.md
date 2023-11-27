# grpcmd

grpcmd is a simple, easy-to-use, and developer-friendly CLI tool for gRPC.

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
The above command will install `grpcmd` package from the tap [`grpcmd/tap`](https://github.com/grpcmd/homebrew-tap).

### Binary
You can also download the binary files for macOS, Linux, and Windows from the [Releases](https://github.com/grpcmd/grpcmd/releases) page.

## Usage

### Start a gRPC Server

    grpc server

Output:

    Listening on address: :50051

    Try running:
            grpc :50051 UnaryMethod '{"name": "Bob"}'

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

### Sending an Empty Request
To send an empty request, simply pass an empty argument for the request data. For example:

    grpc :50051 UnaryMethod ''

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