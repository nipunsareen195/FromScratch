package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"

	"jumbochain.org/p2p"
	"jumbochain.org/temp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourcePort := flag.Int("sp", 0, "source port number")
	flag.Parse()

	fmt.Println("hi")
	fmt.Println("hi")
	fmt.Println("hi")

	r := rand.Reader

	h, err := p2p.MakeHost(*sourcePort, r)
	if err != nil {
		log.Fatal(err)
	}

	if *sourcePort != 0 {
		p2p.StartListener(ctx, h, *sourcePort)
		<-ctx.Done()
	} else {

		peerlist := temp.ReadCsv("peerlist.csv")

		for i := 0; i < len(peerlist); i++ {
			fmt.Println(peerlist[i][0])
			target := peerlist[i][0]
			info := p2p.RunSender(h, target)

			p2p.SendStream(h, info)
			<-ctx.Done()
		}
	}
}
