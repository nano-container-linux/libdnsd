package libdnsd

import (
    "context"
    "net"
    "testing"

    "google.golang.org/grpc"
)

func TestSubmitDynamicOverGRPC_Server(t *testing.T) {
    // start a local grpc server using jsonCodec
    lis, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        t.Fatalf("listen: %v", err)
    }
    srv := grpc.NewServer(grpc.ForceServerCodec(jsonCodec{}))

    // register a minimal service descriptor for dnsd.DynamicDNS with Submit
    sd := &grpc.ServiceDesc{
        ServiceName: "dnsd.DynamicDNS",
        HandlerType: (*interface{})(nil),
        Methods: []grpc.MethodDesc{
            {
                MethodName: "Submit",
                Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
                    var req DynamicSubmitRequest
                    if err := dec(&req); err != nil {
                        return nil, err
                    }
                    return &DynamicSubmitResponse{Message: "received", ID: "abc123", Path: "/tmp"}, nil
                },
            },
        },
    }
    srv.RegisterService(sd, struct{}{})

    go srv.Serve(lis)
    defer srv.Stop()

    target := lis.Addr().String()
    req := DynamicSubmitRequest{PayloadHCL: "x"}
    resp, err := SubmitDynamicOverGRPC(target, req)
    if err != nil {
        t.Fatalf("SubmitDynamicOverGRPC error: %v", err)
    }
    if resp.ID != "abc123" || resp.Message != "received" {
        t.Fatalf("unexpected resp: %+v", resp)
    }
}
