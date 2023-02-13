package p2p

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"

	ma "github.com/multiformats/go-multiaddr"

	"jumbochain.org/temp"
)

func MakeHost(port int, randomness io.Reader) (host.Host, error) {
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 0.0.0.0 will listen on any interface device.
	// ipAddress := GetOutboundIP()

	sourceMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/192.168.1.230/tcp/%d", port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	return libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
}

func StartListener(ctx context.Context, ha host.Host, listenPort int) {
	fullAddr := getHostAddress(ha)
	log.Printf("I am %s\n", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("listener received new stream")
		if err := doEcho1(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	log.Println("listening for connections")

	log.Printf("Now run \"./echo -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)

}

func RunSender(ha host.Host, targetPeer string) peer.AddrInfo {
	fullAddr := getHostAddress(ha)
	log.Printf("I am %s\n", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("sender received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	// Turn the targetPeer into a multiaddr.
	maddr, err := ma.NewMultiaddr(targetPeer)
	if err != nil {
		log.Println(err)
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println(err)
	}

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	log.Println("sender opening stream")

	return *info
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol

	// s, err := ha.NewStream(context.Background(), info.ID, "/echo/1.0.0")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// log.Println("sender saying hello")

	// message := "ye hai message"

	// fmt.Println(message + "\n")

	// _, err = s.Write([]byte(message + "\n"))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// out, err := io.ReadAll(s)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// log.Printf("read reply: %q\n", out)
}

func SendStream(ha host.Host, info peer.AddrInfo) {

	trxs_inMemPool := temp.ReadCsv("TrxMemPool.csv")

	for i := 0; i < len(trxs_inMemPool); i++ {

		trx_details := trxs_inMemPool[i][0]

		s, err := ha.NewStream(context.Background(), info.ID, "/echo/1.0.0")
		if err != nil {
			log.Println(err)
			return
		}

		// log.Println("sender saying hello")

		// message := "ye hai message"

		log.Println(trx_details)

		_, err = s.Write([]byte(trx_details + "\n"))
		fmt.Println(err)
		if err != nil {
			log.Println(err)
			return
		}

		out, err := io.ReadAll(s)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("read reply: %q\n", out)

	}

	// do stuff

}

func SendStream1(ha host.Host, info peer.AddrInfo) {

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				//ee

				trxs_inMemPool := temp.ReadCsv("TrxMemPool.csv")

				for i := 0; i < len(trxs_inMemPool); i++ {

					trx_details := trxs_inMemPool[i][0]

					s, err := ha.NewStream(context.Background(), info.ID, "/echo/1.0.0")
					if err != nil {
						log.Println(err)
						return
					}

					// log.Println("sender saying hello")

					// message := "ye hai message"

					log.Println(trx_details)

					_, err = s.Write([]byte(trx_details + "\n"))
					fmt.Println(err)
					if err != nil {
						log.Println(err)
						return
					}

					out, err := io.ReadAll(s)
					if err != nil {
						log.Println(err)
						return
					}

					log.Printf("read reply: %q\n", out)

					if err := os.Truncate("TrxMemPool.csv", 0); err != nil {
						log.Printf("Failed to truncate: %v", err)
					}

				}

				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Println("0000000")
	log.Printf("read : %s", str)
	_, err = s.Write([]byte(str))
	return err
}

func doEcho1(s network.Stream) error {
	buf := bufio.NewReader(s)

	// buf.ReadBytes('\n')

	str, err := buf.ReadBytes('\n')
	if err != nil {
		return err
	}

	log.Println("111111111")
	log.Printf("read : %s", str)

	temp.UpdateCsv("TrxMemPoolValidator.csv", string(str))

	// byy := temp.DecodeToPerson(str)

	// fmt.Println(byy)

	// _, err = s.Write([]byte(str))
	return err
}
