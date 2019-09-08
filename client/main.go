/*
@Time : 2019/9/8 16:20
@Author : louis
@File : main
@Software: GoLand
*/

package main

import (
	pb "../api"
	"context"
	"flag"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of say backend")
	output := flag.String("o", "output.wav", "output will be written")
	flag.Parse()

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connet to %s: %v", *backend, err)
	}
	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)

	text := &pb.Text{Text: "hello there"}
	res, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatalf("could not say %s: %v", text.Text, err)
	}

	if err = ioutil.WriteFile(*output, res.Audio, 0666); err != nil {
		log.Fatalf("cloud not witte file")
	}
}

