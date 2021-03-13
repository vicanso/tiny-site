package service

import (
	"context"
	"errors"
	"time"

	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/util"
	"google.golang.org/grpc"

	pb "github.com/vicanso/tiny/pb"
)

var (
	grpcConn *grpc.ClientConn
)

type (
	// ImageOptimParams image optim params
	ImageOptimParams struct {
		Data       []byte
		Type       string
		SourceType string
		Quality    int
		Width      int
		Height     int
		Crop       int
	}
	// OptimSrv optim service
	OptimSrv struct{}
)

func init() {
	done := make(chan int)
	go func() {
		opts := make([]grpc.DialOption, 0)
		opts = append(opts, grpc.WithInsecure())
		if util.IsProduction() {
			opts = append(opts, grpc.WithBlock())
		}

		conn, err := grpc.Dial(config.GetTinyAddress(), opts...)
		if err != nil {
			panic(err)
		}
		done <- 1
		grpcConn = conn
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		panic(errors.New("grpc dial timeout"))
	}

}

// Image image optim
func (srv *OptimSrv) Image(params ImageOptimParams) (data []byte, err error) {
	client := pb.NewOptimClient(grpcConn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	in := &pb.OptimRequest{
		Data:    params.Data,
		Quality: uint32(params.Quality),
		Width:   uint32(params.Width),
		Height:  uint32(params.Height),
		Crop:    uint32(params.Crop),
	}
	switch params.Type {
	case "png":
		in.Output = pb.Type_PNG
	case "webp":
		in.Output = pb.Type_WEBP
	case "avif":
		in.Output = pb.Type_AVIF
	default:
		in.Output = pb.Type_JPEG
	}
	switch params.SourceType {
	case "png":
		in.Source = pb.Type_PNG
	case "webp":
		in.Source = pb.Type_WEBP
	default:
		in.Source = pb.Type_JPEG
	}

	reply, err := client.DoOptim(ctx, in)
	if err != nil {
		return
	}
	data = reply.Data
	return
}
