package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"jumbochain.org/p2p"
	"jumbochain.org/temp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourcePort := flag.Int("sp", 0, "source port number")
	flag.Parse()

	r := rand.Reader

	h, err := p2p.MakeHost(*sourcePort, r)
	if err != nil {
		log.Fatal(err)
	}

	if *sourcePort != 0 {
		p2p.StartListener(ctx, h, *sourcePort)
		<-ctx.Done()
	} else {
		ticker := time.NewTicker(10 * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:

					peerlist := temp.ReadCsv("peerlist.csv")

					for i := 0; i < len(peerlist); i++ {
						fmt.Println(peerlist[i][0])
						target := peerlist[i][0]
						info := p2p.RunSender(h, target)

						p2p.SendStream(h, info)

					}
					if err := os.Truncate("TrxMemPool.csv", 0); err != nil {
						log.Printf("Failed to truncate: %v", err)
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
		<-ctx.Done()
	}
}
