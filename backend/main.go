package main

import (
	pb "github.com/oldthreefeng/say-rpc/api"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)

func main() {
	port := flag.Int("p", 8080, "prot to listen")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Printf("could not list to port %d:%v", *port, err)
	}
	cmd := exec.Command("flite", "-t", os.Args[1], "output.wav")
	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	err := s.Server(lis)
	if err != nil {
		log.Fatalf("cloud not server: %v", err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

type server struct {
}

func (server) Say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("cloud not create tmp file : %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("cloud not close tmp file : %v", err)
	}

	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())

	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("file failed : %s", data)
	}
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read tmp file : %s", data)
	}
	return &pb.Speech{Audio: data}, nil
}

