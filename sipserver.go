package main

import (
	"github.com/superirale/sipserver/sip"
)

func main() {
	servicePort := ":5060"
	sip.UseUDP(servicePort)
	// server.UseTCP(servicePort)
}
