package main

import (
	"fmt"
	"log"

	pb "github.com/sempr/ips_go/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIPSVCClient(conn)

	// Contact the server and print out its response.
	for j := 0; j < 256; j += 19 {
		i := j % 256
		r, err := c.IPQuery(context.Background(), &pb.IPRequest{Ip: fmt.Sprintf("%d.%d.%d.%d", i, i, i, i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		_ = r
		log.Printf("One by One : %20s %-30s %s", r.Ip, r.City, r.Loc)
	}
	var rqs []*pb.IPRequest
	for j := 0; j < 256; j += 19 {
		i := j
		rq := pb.IPRequest{Ip: fmt.Sprintf("%d.%d.%d.%d", i, i, i, i)}
		rqs = append(rqs, &rq)
	}
	rq := pb.IPsRequest{Ips: rqs}
	r, err := c.IPSQuery(context.Background(), &rq)
	if err != nil {
		log.Fatalf("could not xxx: %v", err)
	}
	for _, x := range r.Ipr {
		log.Printf("Batch Query: %20s %-30s %s", x.Ip, x.City, x.Loc)
	}

	stream, _ := c.IPStreamQuery(context.Background())

	func() {
		for i := 0; i < 256; i += 19 {
			msg := &pb.IPRequest{Ip: fmt.Sprintf("%d.%d.%d.%d", i, i, i, i)}
			stream.Send(msg)
		}
		stream.CloseSend()
	}()

	func() {

		for {
			r2, err := stream.Recv()
			if r2 == nil {
				fmt.Printf("%#v %#v", r2, err)
				break
			}
			fmt.Println(r2.Ip, r2.City, r2.Loc)
		}
	}()
}
