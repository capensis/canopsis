package bulk

import (
	"fmt"

	influxmod "github.com/influxdata/influxdb/client/v2"
)

type unsafe struct {
	batchCount  int
	batchSize   int
	batchConfig influxmod.BatchPointsConfig
	batch       influxmod.BatchPoints
	client      influxmod.Client
}

func (b *unsafe) reset() error {
	var err error
	b.batchCount = 0
	b.batch, err = influxmod.NewBatchPoints(b.batchConfig)
	if err != nil {
		return fmt.Errorf("influx/Bulk.reset: %v", err)
	}
	return nil
}

func (b *unsafe) AddPoints(po ...PointOp) error {
	b.batchCount += len(po)

	for idx, point := range po {
		ipoint, err := influxmod.NewPoint(
			point.Name,
			point.Tags,
			point.Fields,
			point.Time,
		)
		if err != nil {
			return fmt.Errorf("influx/bulk.AddPoints: point %d: %v", idx, err)
		}
		b.batch.AddPoint(ipoint)
	}

	if b.batchCount > b.batchSize {
		if err := b.Perform(); err != nil {
			return fmt.Errorf("influx/bulk.AddPoints->Perform: %v", err)
		}
	}

	return nil
}

func (b *unsafe) Perform() error {
	if err := b.client.Write(b.batch); err != nil {
		return fmt.Errorf("influx/bulk.Perform: %v", err)
	}
	return b.reset()
}

// New creates a Bulk instance which is not thread safe.
func New(batchSize int, client influxmod.Client, config influxmod.BatchPointsConfig) (Bulk, error) {
	b := unsafe{
		batchConfig: config,
		batchSize:   batchSize,
		client:      client,
	}
	if err := b.reset(); err != nil {
		return nil, fmt.Errorf("influx/Bulk.New: %v", err)
	}

	return &b, nil

}
