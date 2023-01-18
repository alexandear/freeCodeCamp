package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

const (
	scanOpenPortsTimeout = 60 * time.Second

	lookupTimeout = time.Second

	scanPortTimeout = 500 * time.Millisecond
	scanPortNetwork = "tcp"
)

var (
	ErrorInvalidHostname  = errors.New("invalid hostname")
	ErrorInvalidIPAddress = errors.New("invalid IP address")
)

type OpenPortsResult struct {
	ports   []int
	host    string
	address string
}

func OpenPorts(target string, portStart, portEnd int) (OpenPortsResult, error) {
	if target == "" {
		return OpenPortsResult{}, errors.New("target must be non-empty")
	}
	if portStart < 0 || portEnd < 0 {
		return OpenPortsResult{}, errors.New("ports must be positive")
	}
	if portStart > portEnd {
		return OpenPortsResult{}, errors.New("portStart must be lower portEnd")
	}

	ctx, cancel := context.WithTimeout(context.Background(), scanOpenPortsTimeout)
	defer cancel()

	var host, address string
	if ifCanBeIP4(target) {
		address = target
		h, err := lookupHost(ctx, address)
		if err != nil {
			return OpenPortsResult{}, err
		}
		host = h
	} else {
		host = target
		addr, err := lookupAddress(ctx, host)
		if err != nil {
			return OpenPortsResult{}, err
		}
		address = addr
	}

	ports := make([]int, 0, portEnd-portStart)
	for port := portStart; port <= portEnd; port++ {
		if scanPort(ctx, target, port) {
			ports = append(ports, port)
		}
	}

	if ctx.Err() != nil {
		return OpenPortsResult{}, ctx.Err()
	}

	return OpenPortsResult{
		ports:   ports,
		host:    host,
		address: address,
	}, nil
}

func (o *OpenPortsResult) Ports() []int {
	return o.ports
}

func (o *OpenPortsResult) Verbose() string {
	var sb strings.Builder
	sb.Grow(100)
	_, _ = sb.WriteString("Open ports for ")
	if o.host != "" {
		_, _ = sb.WriteString(o.host)
	}
	if o.address != "" && o.host == "" {
		_, _ = sb.WriteString(o.address)
	} else if o.address != "" {
		_, _ = sb.WriteString(" (")
		_, _ = sb.WriteString(o.address)
		_, _ = sb.WriteString(")")
	}
	_, _ = sb.WriteString("\nPORT     SERVICE")
	for _, p := range o.ports {
		port := fmt.Sprintf("\n%-9d", p)
		_, _ = sb.WriteString(port)
		_, _ = sb.WriteString(portsServices[p])
	}

	return sb.String()
}

func lookupHost(ctx context.Context, address string) (string, error) {
	if _, err := netip.ParseAddr(address); err != nil {
		return "", ErrorInvalidIPAddress
	}

	lookupCtx, lookupCancel := context.WithTimeout(ctx, lookupTimeout)
	defer lookupCancel()

	names, err := net.DefaultResolver.LookupAddr(lookupCtx, address)
	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], "."), nil
	}

	return "", nil
}

func lookupAddress(ctx context.Context, host string) (string, error) {
	lookupCtx, lookupCancel := context.WithTimeout(ctx, lookupTimeout)
	defer lookupCancel()

	addresses, err := net.DefaultResolver.LookupHost(lookupCtx, host)
	if err != nil {
		return "", ErrorInvalidHostname
	}

	if len(addresses) > 0 {
		return addresses[0], nil
	}

	return "", nil
}

func scanPort(ctx context.Context, target string, port int) bool {
	scanCtx, scanCancel := context.WithTimeout(ctx, scanPortTimeout)
	defer scanCancel()

	address := target + ":" + strconv.Itoa(port)

	var d net.Dialer
	conn, err := d.DialContext(scanCtx, scanPortNetwork, address)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()

	return true
}

func ifCanBeIP4(target string) bool {
	octets := strings.Split(target, ".")
	if len(octets) != 4 {
		return false
	}
	for _, oct := range octets {
		_, err := strconv.Atoi(oct)
		if err != nil {
			return false
		}
	}
	return true
}
