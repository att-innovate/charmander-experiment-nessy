package main

import (
	"fmt"
	dns "github.com/miekg/dns"
)

func main() {
	message := new(dns.Msg)
	message.SetQuestion("www.google.com.", dns.TypeA)

	client := new(dns.Client)
	response, time, _:= client.Exchange(message, "172.31.2.12:53")
	fmt.Println(response, time)

}
