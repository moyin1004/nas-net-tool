package main

import (
	"log"
	"net"
	"net/netip"

	"github.com/mdlayher/ndp"
)

func GetRa(interfaceName string) (string, error) {
	// Select a network interface by its name to use for NDP communications.
	ifi, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", err
	}

	// Set up an *ndp.Conn, bound to this interface's link-local IPv6 address.
	c, ip, err := ndp.Listen(ifi, ndp.LinkLocal)
	if err != nil {
		return "", err
	}
	// Clean up after the connection is no longer needed.
	defer c.Close()

	log.Println("ndp: bound to address:", ip)
	// ndp: bound to address: fe80::76d4:35ff:fee7:cbc4

	// Choose a target with a known IPv6 link-local address.
	// target := net.ParseIP("fe80::5054:ff:fe33:997")
	target := netip.IPv6LinkLocalAllRouters()
	// Use target's solicited-node multicast address to request that the target
	// respond with a neighbor advertisement.
	_, err = ndp.SolicitedNodeMulticast(target)
	if err != nil {
		return "", err
	}

	// Build a router solicitation message, specifying our source link-layer
	// address so the router does not have to ask it for it explicitly.
	m := &ndp.RouterSolicitation{
		Options: []ndp.Option{
			&ndp.LinkLayerAddress{
				Direction: ndp.Source,
				Addr:      ifi.HardwareAddr,
			},
		},
	}

	// Send to the "IPv6 link-local all routers" multicast group and wait
	// for a response.
	if err := c.WriteTo(m, nil, target); err != nil {
		return "", err
	}
	msg, _, from, err := c.ReadFrom()
	if err != nil {
		return "", err
	}

	// Expect a router advertisement message.
	ra, ok := msg.(*ndp.RouterAdvertisement)
	if !ok {
		return "", err
	}

	// Iterate options and display information.
	log.Printf("ndp: router advertisement from %s:\n", from)
	for _, o := range ra.Options {
		switch o := o.(type) {
		case *ndp.PrefixInformation:
			log.Printf("  - prefix %q: SLAAC: %t\n", o.Prefix, o.AutonomousAddressConfiguration)
			return o.Prefix.String(), nil
		case *ndp.LinkLayerAddress:
			log.Printf("  - link-layer address: %s\n", o.Addr)
		}
	}
	// ndp: router advertisement from fe80::618:d6ff:fea1:ceb7:
	//   - prefix "2600:6c4a:787f:d200::": SLAAC: true
	//   - prefix "fd00::": SLAAC: true
	//   - link-layer address: 04:18:d6:a1:ce:b7
	return "", nil
}
