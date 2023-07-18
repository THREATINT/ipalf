package main

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	var (
		app *cli.App
		err error
	)

	app = &cli.App{
		Name:  "ipalf",
		Usage: "IP Address List Filter, please see github.com/THREATINT/ipalf",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "ipv4",
				Value: false,
				Usage: "Allow IPv4 addresses",
			},
			&cli.BoolFlag{
				Name:  "ipv6",
				Value: false,
				Usage: "Allow IPv6 addresses",
			},
			&cli.BoolFlag{
				Name:  "singleip",
				Value: false,
				Usage: "Allow single ip addresses (including /32 for IPv4 and /128 for IPv6)",
			},
			&cli.BoolFlag{
				Name:  "network",
				Value: false,
				Usage: "Allow networks",
			},
		},

		Action: func(cCtx *cli.Context) error {
			return run(cCtx)
		},
	}

	if err = app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(0xff)
	}
}

func run(cCtx *cli.Context) error {
	var (
		err     error
		scanner *bufio.Scanner
		lines   []string
		line    string
		ip      netip.Addr
		network netip.Prefix
	)

	scanner = bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		line := scanner.Text()

		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	for _, line = range lines {
		var (
			ipAddrTest string
			ipAddr     netip.Addr
		)

		line = strings.TrimSpace(line)

		if cCtx.Bool("ipv6") {
			ipAddrTest = strings.TrimSuffix(line, "/128")
			if ipAddr, err = netip.ParseAddr(ipAddrTest); err == nil {
				if ipAddr.Is6() {
					line = ipAddrTest
				}
			}
		}

		if cCtx.Bool("ipv4") {
			ipAddrTest = strings.TrimSuffix(line, "/32")
			if ipAddr, err = netip.ParseAddr(ipAddrTest); err == nil {
				if ipAddr.Is4() {
					line = ipAddrTest
				}
			}
		}

		if cCtx.Bool("singleip") {
			if ip, err = netip.ParseAddr(line); err == nil {
				if (ip.Is4() && cCtx.Bool("ipv4")) || (ip.Is6() && cCtx.Bool("ipv6")) {
					fmt.Println(line)
					continue
				}
			}
		}

		if cCtx.Bool("network") {
			if network, err = netip.ParsePrefix(line); err == nil {
				if (network.Addr().Is4() && cCtx.Bool("ipv4")) || (network.Addr().Is6() && cCtx.Bool("ipv6")) {
					fmt.Println(line)
				}
			}
		}
	}

	return nil
}
