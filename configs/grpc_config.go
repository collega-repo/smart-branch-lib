package configs

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetConnService(host string, dialOption ...grpc.DialOption) (*grpc.ClientConn, error) {
	dialOption = append(dialOption, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return grpc.Dial(host, dialOption...)
}
