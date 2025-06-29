package core

import (
	"fmt"

	"github.com/CoreKitMDK/corekit-service-configuration/v2/pkg/configuration"
	"github.com/CoreKitMDK/corekit-service-logger/v2/pkg/logger"
	"github.com/CoreKitMDK/corekit-service-metrics/v2/pkg/metrics"
	"github.com/CoreKitMDK/corekit-service-tracing/v2/pkg/tracing"
)

type ICore interface {
	Stop() error
}

type Core struct {
	Logger        logger.IMultiLogger
	Metrics       metrics.IMultiMetrics
	Configuration configuration.IConfiguration
	Tracing       tracing.ITracing
}

func NewCore() (*Core, error) {
	telemetry := &Core{}
	var err error = nil

	{
		config := metrics.NewConfiguration()
		config.UseConsole = true
		config.UseNATS = true
		config.NatsURL = "internal-metrics-broker-nats-client"
		config.NatsPassword = "internal-metrics-broker"
		config.NatsUsername = "internal-metrics-broker"
		telemetry.Metrics = config.Init()
	}

	{
		config := logger.NewConfiguration()
		config.UseConsole = true
		config.UseNATS = true
		config.NatsURL = "nats://internal-logger-broker-nats:4222"
		config.NatsPassword = "internal-logger-broker"
		config.NatsUsername = "internal-logger-broker"
		telemetry.Logger = config.Init()
	}

	{
		config := tracing.NewConfiguration()
		config.UseConsole = true
		config.UseNATS = true
		config.NatsURL = "nats://internal-tracing-broker-nats:4222"
		config.NatsPassword = "internal-tracing-broker"
		config.NatsUsername = "internal-tracing-broker"
		telemetry.Tracing = config.Init()
	}

	{
		config := configuration.NewConfiguration()
		config.UseConfigFile = false
		config.UseEnv = false
		config.UseConfigString = false
		config.UseApi = true
		config.ApiUrl = "http://config"
		config.ApiNamespace = "testing-dev" //should be overwritten in API client
		client := config.Init()
		if client == nil {
			err = fmt.Errorf("configuration client is nil")
		}
		telemetry.Configuration = client
	}

	return telemetry, err
}

func (t *Core) Stop() error {
	if t.Metrics != nil {
		t.Metrics.Stop()
	}
	if t.Logger != nil {
		t.Logger.Stop()
	}
	return nil
}
