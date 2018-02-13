package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

const (
	db       = "ggg"
	username = "root"
	password = "root"
	meas     = "cpu_usage"
)

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

func runQuery(c client.Client, cmd string) {
	q := client.NewQuery(cmd, db, "")
	if res, err := c.Query(q); err != nil && res.Error() == nil {
		fmt.Errorf("runQuery() error: %v: %s\n", res.Results, err)
	}
}

func writeData(c client.Client, p Point) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "",
	})
	if err != nil {
		log.Fatal(err)
	}

	pt, err := client.NewPoint(p.name, p.tags, p.fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func fromFile(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("%s could not read: %v\n", filePath, err)
	}
	defer f.Close()

	// read each line by Scanner()
	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "%s scanner error: %v\n", filePath, err)
	}

	return lines
}

func main() {

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// create database if not exists
	runQuery(c, fmt.Sprintf("CREATE DATABASE %s", db))

	// read json each line
	lines := fromFile("a.json")
	d := new(Data)

	for _, v := range lines {
		jsonBytes := ([]byte)(v)
		if err := json.Unmarshal(jsonBytes, d); err != nil {
			fmt.Errorf("JSON Unmarshal error: %s", err)
			return
		}

		// insert data
		p := Point{
			name: meas,
			tags: map[string]string{
				"cpu": "cpu-total",
			},
			fields: map[string]interface{}{
				"server": d.Server,
				"client": d.Client,
				"value":  1,
			},
		}
		writeData(c, p)

	}
}
