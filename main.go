package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	fs := flag.NewFlagSet("persist", flag.ContinueOnError)
	domain := fs.String("d", "", "specify the domain to lookup")
	numwork := fs.Int("n", 300, "specify the number of workers to use")
	fs.Parse(os.Args[1:])
	Main(*domain, *numwork)
}

func Main(d string, n int) {
	log.SetFlags(0)

	ipch := make(chan string)
	go func() {
		if err := producer(ipch); err != nil {
			log.Fatalf("input error: %v", err)
		}
		close(ipch)
	}()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			worker(ipch, d)
			wg.Done()
		}()
	}
	wg.Wait()
}

func producer(ipch chan string) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ipch <- strings.TrimSpace(scanner.Text())
	}
	return scanner.Err()
}

func worker(ipch chan string, domain string) {
	for ip := range ipch {
		if err := lookupMX(domain, ip); err != nil {
			log.Printf("%v <- DNS server || error -> %v", ip, err)
		}
	}
}

func lookupMX(name, dnsIP string) error {
	if strings.ContainsRune(dnsIP, ':') { // ipv6
		dnsIP = "[" + dnsIP + "]"
	}
	dns := net.Resolver{
		PreferGo:     true,
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, network, dnsIP+":53")
		},
	}

	_, err := dns.LookupMX(context.Background(), name)
	return err
}
