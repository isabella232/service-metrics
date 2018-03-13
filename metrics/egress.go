package metrics

import (
	"errors"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/lager"
)

type Logger interface {
	Info(string, ...lager.Data)
	Error(action string, err error, data ...lager.Data)
}

type GaugeIngressClient interface {
	EmitGauge(opts ...loggregator.EmitGaugeOption)
}

type EgressClient struct {
	emitter    GaugeIngressClient
	sourceID   string
	instanceID int
}

func NewEgressClient(inClient GaugeIngressClient, sourceID string) *EgressClient {
	return &EgressClient{
		emitter:  inClient,
		sourceID: sourceID,
	}
}

func (c *EgressClient) Emit(metrics Metrics, logger Logger) {
	if len(metrics) < 1 {
		logger.Info("sending-metrics", lager.Data{
			"details": "no metrics to send",
		})
		return
	}

	if c.sourceID == "" {
		e := errors.New("You must set a source ID")
		logger.Error("sending metrics failed", e, lager.Data{
			"Emit": "failed",
		})
	}

	logger.Info("sending-metrics", lager.Data{
		"details": "emitting gauges to logging platform",
		"count":   len(metrics),
	})

	opts := []loggregator.EmitGaugeOption{
		loggregator.WithGaugeAppInfo(c.sourceID, c.instanceID),
	}

	for _, m := range metrics {
		opts = append(opts, loggregator.WithGaugeValue(m.Key, m.Value, m.Unit))
	}

	c.emitter.EmitGauge(opts...)
}

func (c *EgressClient) SetInstanceID(instanceID int) {
	c.instanceID = instanceID
}
