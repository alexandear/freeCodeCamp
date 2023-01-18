package main

import (
	"fmt"
	"log"
)

func main() {
	// Called with URL
	res, err := OpenPorts("www.freecodecamp.org", 75, 85)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Open ports:", res.Ports())

	// Called with ip address
	res, err = OpenPorts("104.26.10.78", 8079, 8090)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Open ports:", res.Ports())

	// Verbose called with ip address and no host name returned -- single open port
	res, err = OpenPorts("104.26.10.78", 440, 450)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Verbose())

	// Verbose called with ip address and valid host name returned -- single open port
	res, err = OpenPorts("137.74.187.104", 440, 450)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Verbose())

	// Verbose called with host name -- multiple ports returned
	res, err = OpenPorts("scanme.nmap.org", 20, 80)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Verbose())
}
