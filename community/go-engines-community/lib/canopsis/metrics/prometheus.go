package metrics

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"time"
)

const Namespace = "metrics"

const (
	DurationHour  = time.Hour
	DurationDay   = 24 * DurationHour
	DurationWeek  = 7 * DurationDay
	DurationMonth = 30 * DurationDay
)

type MetricPoint struct {
	Timestamp int64 `json:"timestamp"`
	Value     int64 `json:"value"`
}

type prometheusAPIClient struct {
	prometheusApi v1.API
	logger        zerolog.Logger
}

type PrometheusAPIClient interface {
	/**
	For example: to get ack alarms number with sampling by hours

	res, err := a.prometheusClient.GetCountersVector(
		c.Request.Context(),
		"metrics_ack_alarm_number_slice",
		"slice=\"default_slice\"",
		time.Now().Add(-time.Hour * 4),
		time.Now().Add(time.Hour * 2),
		metrics.DurationHour,
	)
	*/
	GetCountersVector(ctx context.Context, metric, labels string, start, end time.Time, stepDuration time.Duration) ([]MetricPoint, error)
	/**
	For example: to get average resolve time with sampling by hours

	res, err := a.prometheusClient.GetAverageDurationVector(
		c.Request.Context(),
		"metrics_resolve_duration_slice_sum",
		"metrics_resolve_duration_slice_count",
		"slice=\"default_slice\"",
		time.Now().Add(-time.Hour * 4),
		time.Now().Add(time.Hour * 2),
		metrics.DurationHour,
	)
	*/
	GetAverageDurationVector(ctx context.Context, metricSum, metricCount, labels string, start, end time.Time, stepDuration time.Duration) ([]MetricPoint, error) // only for histograms or summary
	/**

	For example: to get percentage of ack alarms in comparison to the total number with sampling by hours

	res, err := a.prometheusClient.GetRatiosVector(
		c.Request.Context(),
		"metrics_ack_alarm_number_slice",
		"metrics_total_alarm_number_slice",
		"slice=\"default_slice\"",
		time.Now().Add(-time.Hour * 4),
		time.Now().Add(time.Hour * 2),
		metrics.DurationHour,
	)
	*/
	GetRatiosVector(ctx context.Context, metricMain, metricTotal, labels string, start, end time.Time, stepDuration time.Duration) ([]MetricPoint, error)
}

func NewPrometheusAPIClient(
	prometheusAddress string,
	logger zerolog.Logger,
) (PrometheusAPIClient, error) {
	client, err := api.NewClient(api.Config{
		Address: prometheusAddress,
	})
	if err != nil {
		return nil, err
	}

	return &prometheusAPIClient{
		prometheusApi: v1.NewAPI(client),
		logger:        logger,
	}, nil
}

func (p *prometheusAPIClient) GetAverageDurationVector(ctx context.Context, metricSum, metricCount, labels string, start, end time.Time, stepDuration time.Duration) ([]MetricPoint, error) {
	return p.getVectorResult(
		ctx,
		fmt.Sprintf("sum(rate(%s{%s}[%dh]))/sum(rate(%s{%s}[%dh]))", metricSum, labels, int(stepDuration.Hours()), metricCount, labels, int(stepDuration.Hours())),
		start,
		end,
		stepDuration,
	)
}

func (p *prometheusAPIClient) GetCountersVector(ctx context.Context, metric, labels string, start, end time.Time, stepDuration time.Duration) ([]MetricPoint, error) {
	return p.getVectorResult(
		ctx,
		fmt.Sprintf("sum(increase(%s{%s}[%dh]))", metric, labels, int(stepDuration.Hours())),
		start,
		end,
		stepDuration,
	)
}

func (p *prometheusAPIClient) GetRatiosVector(ctx context.Context, metricMain, metricTotal, labels string, start, end time.Time, stepDuration time.Duration) ([]MetricPoint, error) {
	return p.getVectorResult(
		ctx,
		fmt.Sprintf("sum(increase(%s{%s}[%dh]))/sum(increase(%s{%s}[%dh]))*100", metricMain, labels, int(stepDuration.Hours()), metricTotal, labels, int(stepDuration.Hours())),
		start,
		end,
		stepDuration,
	)
}

func (p *prometheusAPIClient) getVectorResult(ctx context.Context, query string, start, end time.Time, stepDuration time.Duration) ([]MetricPoint, error) {
	start = start.Truncate(stepDuration)
	end = end.Truncate(stepDuration)

	queryResult, warnings, err := p.prometheusApi.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  stepDuration,
	})
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return nil, err
	}

	result, ok := queryResult.(model.Matrix)
	if !ok {
		return nil, errors.New("result should be prometheus Matrix")
	}

	if len(result) == 0 {
		return nil, errors.New("empty prometheus result")
	}

	for _, wrn := range warnings {
		p.logger.Warn().Msg(wrn)
	}

	stepInSeconds := int64(stepDuration.Seconds())
	previousTimeStamp := start.Add(-time.Duration(stepInSeconds) * time.Second).Unix()
	vector := make([]MetricPoint, 0)

	for _, e := range result[0].Values {
		timestamp := e.Timestamp.Time().Unix()
		if timestamp-previousTimeStamp > stepInSeconds {
			//prometheus doesn't return value pairs for a time period if metric = 0, so just fill the gaps with 0
			for i := previousTimeStamp + stepInSeconds; i < timestamp; i += stepInSeconds {
				vector = append(vector, MetricPoint{
					Timestamp: i,
					Value:     0,
				})
			}
		}

		vector = append(vector, MetricPoint{
			Timestamp: timestamp,
			Value:     int64(e.Value),
		})

		previousTimeStamp = timestamp
	}

	// if there is still some place left, fill it with 0
	if vector[len(vector)-1].Timestamp < end.Unix() {
		for i := vector[len(vector)-1].Timestamp + stepInSeconds; i <= end.Unix(); i += stepInSeconds {
			vector = append(vector, MetricPoint{
				Timestamp: i,
				Value:     0,
			})
		}
	}

	return vector, nil
}
