package service

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	pb "github.com/vicanso/tiny/proto"
)

var (
	grpcConn            *grpc.ClientConn
	defaultOptimTimeout = 10 * time.Second
)

type (
	// OptimOptions optim options
	OptimOptions struct {
		Type      string
		ClipType  string
		ImageType string
		Width     int
		Height    int
		Quality   int
		Data      []byte
	}
)

func init() {
	address := viper.GetString("tiny.address")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	grpcConn = conn
	timeout := viper.GetDuration("tiny.timeout")
	if timeout != 0 {
		defaultOptimTimeout = timeout
	}
}

// GetClipType get clip type
func GetClipType(str string) pb.ClipType {
	switch str {
	case "center":
		return pb.ClipType_CENTER
	case "lt":
		return pb.ClipType_LT
	case "tc":
		return pb.ClipType_TC
	default:
		return pb.ClipType_NONE
	}
}

func getCompressRequest(opts *OptimOptions) (in *pb.CompressRequest) {
	in = &pb.CompressRequest{
		Quality: uint32(opts.Quality),
		Width:   uint32(opts.Width),
		Height:  uint32(opts.Height),
		Data:    opts.Data,
	}
	t := strings.ToUpper(opts.Type)
	in.Type = pb.Type(pb.Type_value[t])
	t = strings.ToUpper(opts.ImageType)
	in.ImageType = pb.Type(pb.Type_value[t])
	if opts.ClipType != "" {
		in.ClipType = GetClipType(opts.ClipType)
	}
	return
}

// Optim optim image
func Optim(opts *OptimOptions) (data []byte, err error) {
	c := pb.NewCompressClient(grpcConn)
	ctx, cancel := context.WithTimeout(context.Background(), defaultOptimTimeout)
	defer cancel()
	r, err := c.Do(ctx, getCompressRequest(opts))

	if err != nil {
		return
	}
	data = r.Data
	return
}
