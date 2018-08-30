package main

import (
        "flag"
        "fmt"
        "log"
        "net"

        metrics "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v2"
        "google.golang.org/grpc"
)

type MetricsService struct{}

func (srv *MetricsService) StreamMetrics(stream metrics.MetricsService_StreamMetricsServer) error {
        for {
                msg, err := stream.Recv()
                if err != nil || msg == nil {
                        return err
                }

                for _, mf := range msg.GetEnvoyMetrics() {
                        for _, m := range mf.GetMetric() {
                                fmt.Printf("%#v\n", m)
                        }
                }
        }
}

func main() {
        listen := flag.String("listen", "127.0.0.1:4994", "Host:port of metrics service")
        flag.Parse()

        tcpListener, err := net.Listen("tcp", *listen)
        if err != nil {
                log.Fatalf("listen failed: %s", err.Error())
                return
        }

        s := grpc.NewServer()
        metrics.RegisterMetricsServiceServer(s, &MetricsService{})

        if err := s.Serve(tcpListener); err != nil {
                log.Fatalf("serve failed: %s", err.Error())
        }
}
