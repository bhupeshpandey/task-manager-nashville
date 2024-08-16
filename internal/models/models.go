package models

type Config struct {
	GrpcServer *GRPCServer `yaml:"grpcServer"`
}

type GRPCServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
