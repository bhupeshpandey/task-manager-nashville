package models

type Config struct {
	GrpcServer    *GRPCServer    `yaml:"grpcServer"`
	MetricsConfig *MetricsConfig `yaml:"prometheus"`
	LoggingConfig *LoggingConfig `yaml:"logging"`
	HttpServer    *HttpServer    `yaml:"httpServer"`
}

type GRPCServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type HttpServer struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Logger Logger
}

type MetricsConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Logger Logger
}

// LoggingConfig represents the logging-related configuration
type LoggingConfig struct {
	Type        string `yaml:"type"`
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"logLevel"`
}
