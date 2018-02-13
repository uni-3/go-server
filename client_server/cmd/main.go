package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	endpoint = flag.String("endpoint", "http://localhost:8080", "endpoint for get request")
)

func main() {
	flag.Parse()
	fmt.Printf("endpoint: %s \n", *endpoint)

	queue := make(chan []byte, 1) // size 1のqueue

	go Goget(*endpoint, queue)
	go Gostore(queue)

	select {}
}

func Get(endpoint string) ([]byte, error) {
	//queryにenvのHOSTNAME=72d6f244cb7dを送信する
	host := os.Getenv("HOSTNAME")
	res, err := http.Get(endpoint + "?host=" + host)
	if err != nil {
		fmt.Errorf("got error http.Get(): %s", err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("got error read body: %s", err)
		return nil, err
	}

	return resBody, nil
}

func Goget(endpoint string, queue chan []byte) {

	for {
		data, err := Get(endpoint)
		if err != nil {
			fmt.Println("got error http.Get(): %s", err)
		}
		queue <- data
		log.Println("get data!!")
		time.Sleep(1 * time.Second)
	}

}

func Store(data []byte) {
	// write file (0777)
	logfile, err := os.OpenFile("./logs/getter.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open .log file:" + err.Error())
	}
	defer logfile.Close()

	// io.MultiWriteで、
	// 標準出力とファイルの両方を束ねて、
	// logの出力先に設定する
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(string(data))
}

func Gostore(queue chan []byte) {
	for {
		data := <-queue
		Store(data)
	}
}
