package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fopina/cfdd/cfapi"
	flag "github.com/spf13/pflag"
)

func getIP() ([]net.IPAddr, error) {
	var resolver *net.Resolver

	resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("udp", "resolver1.opendns.com:53")
		},
	}

	return resolver.LookupIPAddr(context.Background(), "myip.opendns.com")
}

var version string = "DEV"
var date string

type clioptions struct {
	Domain       string
	Record       string
	Email        string
	Token        string
	PollInterval int
}

func main() {
	options := clioptions{}
	versionPtr := flag.BoolP("version", "v", false, "display version")
	flag.StringVarP(&options.Domain, "domain", "d", "", "Domain (or Zone) that the record is part of")
	flag.StringVarP(&options.Record, "record", "r", "", "Record to be updated")
	flag.StringVarP(&options.Email, "email", "e", "", "Email for authentication")
	flag.StringVarP(&options.Token, "token", "t", "", "API token for authentication")
	flag.IntVarP(&options.PollInterval, "polling", "p", 0, "Number of seconds between each check for external IP (use 0 to run only once)")
	helpPtr := flag.BoolP("help", "h", false, "this")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		return
	}

	if *versionPtr {
		fmt.Println("Version: " + version + " (built on " + date + ")")
		return
	}

	client := cfapi.NewCFClient(options.Email, options.Token)

	zone, err := client.FindZoneByName(options.Domain)
	if err != nil {
		log.Fatal(err)
	}

	record, err := client.FindRecordByName(zone.ID, options.Record)

	if record.Type != "A" {
		log.Fatalf("%v does not have a valid record type", record)
	}

	for {
		ips, err := getIP()
		if err != nil {
			log.Printf("failed to retrieve current IP - %v\n", err)
		} else if ips[0].IP.String() != record.Content {
			fmt.Println("Updating IP...")
			newRecord := record
			newRecord.Content = ips[0].IP.String()
			err = client.UpdateRecord(zone.ID, newRecord)
			if err != nil {
				log.Printf("failed to retrieve current IP - %v\n", err)
			} else {
				record.Content = newRecord.Content
			}
		}
		if options.PollInterval == 0 {
			break
		}
		time.Sleep(time.Duration(options.PollInterval) * time.Second)
	}
}
