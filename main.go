package main

import (
	"flag"
	"fmt"
	"github.com/dmitryrn/resource-md5/config"
	"github.com/dmitryrn/resource-md5/transport"
	"github.com/dmitryrn/resource-md5/worker"
	"os"
)

func main() {
	cfg := config.Config{
		Parallel: 10,
	}
	ParseFlags(&cfg)

	httpClient := transport.NewHTTPClient()

	wrkr := worker.NewWorker(cfg, httpClient)

	ch, errCh := wrkr.Do()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				fmt.Printf("error occurred: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		case pair := <-ch:
			fmt.Println(pair[0], pair[1])
		}
	}
}

func ParseFlags(config *config.Config) {
	flag.Uint64Var(&config.Parallel, "parallel", config.Parallel, "")

	flag.Parse()

	args := flag.Args()

	if config.Parallel < 1 {
		fmt.Println("the 'parallel' argument can't be 0")
		os.Exit(1)
	}

	config.URLs = args
}
