package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	ApplicationName = "http-multiplexer"
	Namespace       = "test-tasks"
)

type (
	Logger struct {
		Level string
	}

	K8S struct {
		Port           uint
		HealthEndpoint string
		MetricEndpoint string
	}

	Multiplexer struct {
		AllRequestsTimeOut   time.Duration
		MaxOneTimeRequests   int
		RequestTimeOut       time.Duration
		MaxUrlsInRequest     int
		MaxActiveConnections int
	}

	Root struct {
		Logger      Logger
		K8S         K8S
		Multiplexer Multiplexer
	}
)

func MustConfigure() *Root {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	readConfErr := viper.ReadInConfig()
	if readConfErr != nil {
		panic(readConfErr)
	}

	cfg := Root{}
	unmarshallErr := viper.Unmarshal(&cfg)
	if unmarshallErr != nil {
		panic(unmarshallErr)
	}

	return &cfg
}
