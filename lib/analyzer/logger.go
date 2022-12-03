package analyzer

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func CreateLog(measurement string, tData map[string]string, fData map[string]interface{}) error {

	_, w, err := GetLogger()
	if err != nil {
		return err
	}
	p := influxdb2.NewPointWithMeasurement(measurement)

	//log add tag value
	for t, v := range tData {
		p = p.AddTag(t, v)
	}

	//log add field value
	for f, v := range fData {
		p = p.AddField(f, v)
	}

	err = w.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}
