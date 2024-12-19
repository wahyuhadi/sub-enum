package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sub/model"
	"sub/services"
	"time"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

var (
	target = flag.String("t", "", "target host")
	key    = flag.String("k", "", "key")
)

func main() {
	flag.Parse()

	if *target == "" {
		log.Println("target cannot be nill")
		os.Exit(1)
	}

	if *key == "" {
		*key = "0000-0000-0000-00000"
	}

	subfinderOpts := &runner.Options{
		OutputFile:         "test.txt",
		JSON:               true,
		Silent:             true,
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		// ResultCallback: func(s *resolve.HostEntry) {
		// callback function executed after each unique subdomain is found
		// },
		// ProviderConfig: "your_provider_config.yaml",
		// and other config related options
	}

	// disable timestamps in logs / configure logger
	log.SetFlags(0)

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		log.Fatalf("failed to create subfinder runner: %v", err)
	}

	output := &bytes.Buffer{}
	// To run subdomain enumeration on a single domain
	if err = subfinder.EnumerateSingleDomain(*target, []io.Writer{output}); err != nil {
		log.Fatalf("failed to enumerate single domain: %v", err)
	}

	prefix := strings.Split(output.String(), "\n")
	var metadata model.SubDomainMetaData
	for _, data := range prefix {
		json.Unmarshal([]byte(data), &metadata)
		metadata.Source = *key
		if scan_active_domain(metadata.Host) {
			services.Els(metadata)
			// fmt.Println(metadata)
		}
		// fmt.Println(metadata.Host, metadata.Source)
	}

	// print the output
}

func scan_active_domain(domain string) bool {
	//
	timeout := 5 * time.Second
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", domain, "80"), timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	return true
}
