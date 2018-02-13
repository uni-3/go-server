package store

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

// 接続先はconstractorで設定
type Saver interface {
	IsExists() (bool, error)
	Store(Data) error
}

type Point struct {
	name   string
	tags   map[string]string
	fields map[string]interface{}
}

type Data struct {
	//Host     string        `json:"host"`
	Client string `json:"url"`
	//Query    []interface{} `json:"query"`
	//Body     string        `json:"body"`
	Server string `json:"hostname"`
}

type Influx struct {
	db     string
	meas   string
	tags   map[string]string
	client client.Client
}

func NewInflux(db string, meas string, tags map[string]string, client client.Client) *Influx {
	c := Influx{
		db:     db,
		meas:   meas,
		tags:   tags,
		client: client,
	}
	return &c
}

func (i Influx) RunQuery(c client.Client, cmd string) {
	q := client.NewQuery(cmd, i.db, "")
	if res, err := i.client.Query(q); err != nil && res.Error() == nil {
		fmt.Println(res.Results)
	}
}

func (i Influx) Store(data Data) error {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.db,
		Precision: "",
	})
	if err != nil {
		log.Fatal(err)
	}

	p := Point{
		name: i.meas,
		tags: i.tags,
		fields: map[string]interface{}{
			"server": data.Server,
			"client": data.Client,
			"value":  1,
		},
	}

	pt, err := client.NewPoint(p.name, p.tags, p.fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := i.client.Write(bp); err != nil {
		log.Fatal(err)
	}

	return nil
}

func IsExists() (bool, error) {

	return true, nil

}
