package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	api "github.com/synerex/synerex_api"
	sxutil "github.com/synerex/synerex_sxutil"
	test "local.packages/proto_test"
	pbase "local.packages/synerex_proto"
)

var (
	nodesrv         = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	port            = flag.Int("port", 10080, "Test Provider Listening Port")
	mu              sync.Mutex
	sxServerAddress string
)

func supplyTestCallback(clt *sxutil.SXServiceClient, sp *api.Supply) {
	tst := &test.Hello{}
	err := proto.Unmarshal(sp.Cdata.Entity, tst)
	if err == nil {
		msg := tst.message
		log.Print(msg)
	} else {
		log.Printf("Err UnMarshal %v", err)
	}
}

func reconnectClient(client *sxutil.SXServiceClient) {
	mu.Lock()
	if client.Client != nil {
		client.Client = nil
		log.Printf("Client reset \n")
	}
	mu.Unlock()
	time.Sleep(5 * time.Second) // wait 5 seconds to reconnect
	mu.Lock()
	if client.Client == nil {
		newClt := sxutil.GrpcConnectServer(sxServerAddress)
		if newClt != nil {
			log.Printf("Reconnect server [%s]\n", sxServerAddress)
			client.Client = newClt
		}
	} else { // someone may connect!
		log.Printf("Use reconnected server\n", sxServerAddress)
	}
	mu.Unlock()
}

func subscribeTestSupply(client *sxutil.SXServiceClient) {
	ctx := context.Background()
	for {
		client.SubscribeSupply(ctx, supplyTestCallback)
		log.Print("Error on subscribe")
		reconnectClient(client)
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

	go subscribeTestSupply(pc_client)

	wg.Wait()

	log.Print("sucsess!!!")
}
