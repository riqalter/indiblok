package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

type stringSliceFlag []string

func (f *stringSliceFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *stringSliceFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// resolveDNS dia bakalan ngecek ip address yang kita kasih ke dns resolver nasional/indihome
func resolveDNS(hostname, dnsServer string) (string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := &net.Dialer{}
			return dialer.DialContext(ctx, "udp", dnsServer)
		},
	}

	// ngeresolve ip hostname
	ips, err := resolver.LookupIPAddr(context.Background(), hostname)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IPs found for %s", hostname)
	}

	// return ip sebagai string
	return ips[0].IP.String(), nil
}

func main() {
	var domainFlags stringSliceFlag

	// cuma cara make flag nya
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [domain1] [domain2] ...\n", os.Args[0])
		flag.PrintDefaults()
	}

	// custom flag
	flag.Var(&domainFlags, "d", "Domain to resolve")

	// parse flags
	flag.Parse()

	domains := flag.Args()

	// kalo make flag -d, masukin semua domain
	if len(domainFlags) > 0 {
		domains = append(domains, domainFlags...)
	}

	// kalo gada domain yang di masukan, keluar.
	if len(domains) == 0 {
		flag.Usage()
		return
	}

	// make dns nasional telkom
	dnsServer := "203.130.196.6:53" // ns1.telkom.net.id dari https://portal-uang.com/indihome/dns/, 53 itu port standar dns

	// ngeloop dan ngecek
	for _, rawURL := range domains {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
			continue
		}

		// kalo dia responnya kosong anggap aja dia https 
		if parsedURL.Scheme == "" {
			rawURL = "https://" + rawURL
			parsedURL, err = url.Parse(rawURL)
			if err != nil {
				fmt.Printf("Error parsing URL: %v\n", err)
				continue
			}
		}

		ipAddress, err := resolveDNS(parsedURL.Hostname(), dnsServer)
		if err != nil {
			fmt.Printf("Error resolving ip address for %s: %v\n", rawURL, err)
			continue
		}

		status := "OK"
		if ipAddress == "36.86.63.185" { // https://www.coderstool.com/cs/mHh31k atau https://ip-api.com/internetpositif.id
			status = "(┬┬﹏┬┬) Blocked"
		}

		fmt.Printf("Domain: %s | Status: %s\n", parsedURL.Hostname(), status)
	}

}
