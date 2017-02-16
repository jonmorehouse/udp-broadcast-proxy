package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type options struct {
	upstream []url.URL
	listen   url.URL
	logger   log.Logger
}

func newOpts(rawUpstreamAddrs, listenAddr string, listenPort int) (*options, error) {
	// build out upstream urls
	rawAddrs := strings.Split(rawUpstreamAddrs, ",")
	upstreams := make(*url.URL, len(rawAddrs))
	for idx, rawAddr := range rawAddrs {
		parsed, err := url.Parse(rawAddr)
		if err != nil {
			return nil, fmt.Errorf("Invalid upstream in -upstreams upstream:%s err:%s", rawAddr, err.Error())
		}
		upstreams[idx] = parsed
	}

	// build out listen url
	listenAddr += strconv.Itoa(listenPort)
	listen, err := url.Parse(listenAddr)
	if err != nil {
		return nil, fmt.Errorf("Invalid listen addr: %s err:%s", listenAddr, err.Error())
	}

	// build out logger
	var buf bytes.Buffer
	logger := log.New(&buf, "udp-broadcast-proxy: ", log.Lshortfile)

	return &options{
		upstreams: upstreams,
		listen:    listen,
		logger:    logger,
	}, nil
}

func main() {
	rawUpstreamAddrs := flag.String("upstreams", "", "comma delimited list of upstreams to broadcast too")
	listenPort := flag.Int("listen-port", 0, "listen port")
	listenAddr := flag.String("listen-addr", "", "listen address")

	flag.Parse()

	opts, err := newOpts(*rawUpstreamAddrs, *listenAddr, *listenPort)
	if err != nil {
		log.Fatal(err)
	}

	server, err := newServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)

	// handle any errors that happen either in startup or shutdown
	go func() {
		if err := <-errCh; err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	// start the server in its own goro
	go func() {
		if err := server.start(); err != nil {
			errCh <- err
		}
	}()

	// wait for signals and attempt to stop the server
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		errCh <- server.stop()
	}()

	// wait until the server either shut down properly
	wg.Wait()
}
