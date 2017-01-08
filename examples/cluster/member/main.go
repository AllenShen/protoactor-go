package main

import (
	"fmt"
	"log"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/consul_cluster"
	"github.com/AsynkronIT/protoactor-go/examples/cluster/shared"
)

const (
	timeout = 1 * time.Second
)

func main() {
	cp, err := consul_cluster.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.StartWithProvider("mycluster", "127.0.0.1:8081", cp)

	sync()
	async()

	console.ReadLine()
}

func sync() {
	hello := shared.GetHelloGrain("abc")
	options := []cluster.GrainCallOption{cluster.WithTimeout(5 * time.Second), cluster.WithRetry(5)}
	res, err := hello.SayHello(&shared.HelloRequest{Name: "GAM"}, options...)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message from SayHello: %v", res.Message)
	log.Println("Starting")
	for i := 0; i < 10000; i++ {
		x := shared.GetHelloGrain(fmt.Sprintf("hello%v", i))
		x.SayHello(&shared.HelloRequest{Name: "GAM"})
	}
	log.Println("Done")
}

func async() {
	hello := shared.GetHelloGrain("abc")
	c, e := hello.AddChan(&shared.AddRequest{A: 123, B: 456})

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			log.Println("Tick..") //this might not happen if res returns fast enough
		case err := <-e:
			log.Fatal(err)
		case res := <-c:
			log.Printf("Result is %v", res.Result)
			return
		}
	}
}
