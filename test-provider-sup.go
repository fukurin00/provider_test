package main

import (
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_api"
	sxutil "github.com/synerex/synerex_sxutil"
	"google.golang.org/protobuf/proto"
	test "local.packages/proto_test"
	pbase "local.packages/synerex_proto"
)

var (
	nodesrv         = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port            = flag.Int("port", 10080, "Test Provider Listening Port")
	mu              sync.Mutex
	sxServerAddress string
)

func TestSupply(msg string, client *sxutil.SXServiceClient) {
	p := &test.Hello{
		Message: msg,
	}

	out, err := proto.Marshal(p)
	if err == nil {
		cont := pb.Content{Entity: out}
		smo := sxutil.SupplyOpts{
			Name:  "test supply",
			Cdata: &cont,
		}

		client.NotifySupply(&smo)
	}
}

func main() {
	flag.Parse()

	channelTypes := []uint32{pbase.TEST_SVC}
	sxServerAddress, rerr := sxutil.RegisterNode(*nodesrv, "TestProvider", channelTypes, nil)
	if rerr != nil {
		log.Fatal("Can't register node ", rerr)
	}
	log.Printf("Connecting SynerexServer at [%s]\n", sxServerAddress)

	wg := sync.WaitGroup{}

	client := sxutil.GrpcConnectServer(sxServerAddress) // if there is server address change, we should do it!

	//wg.Add(1)
	if client == nil {
		log.Fatal("Can't connect Synerex Server")
	} else {
		log.Print("Connecting SynerexServer")
	}

	pc_client := sxutil.NewSXServiceClient(client, pbase.TEST_SVC, "{Client:test-provider}")

	wg.Add(1)

	//go subscribeTestSupply(pc_client)
	go TestSupply("hellloooooooooo", pc_client)

	wg.Wait()

	log.Print("sucsess!!!")
}
