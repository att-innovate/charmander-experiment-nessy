package main


import (
	"fmt"
	"time"
	"math/rand"
	"github.com/miekg/dns"
	"os"
	"bufio"
	"sync"
	"flag"
)


type query struct {
	ip string
	dnstype uint16
}



var fileop = flag.String("d", "tiny_20", "datafile")
var numUsers = flag.Int("u",10,"number of users(threads) sending queries")



func main() {
	flag.Parse()
	fmt.Println("the flag is : ",*fileop)
	fmt.Println("the number of users is : ",*numUsers)




	fmt.Println("Collecting queries from ",*fileop, ".........")

	queries := collectQueries(*fileop)
	numQueries := len(queries)
	fmt.Println("Done collecting ",numQueries, " queries")

	var wg sync.WaitGroup
	results := make(chan string, len(queries))

	//random delay with distribution set up
	r := rand.New(rand.NewSource(99))
	StdDev := 300.0
	Mean := 500.0
	fmt.Println("Start sending queries from multiple threads, with delay time normally distributed with Mean:",Mean, " and StdDev:",StdDev)

	for _,query := range queries {
		wg.Add(1)
		go func (ip string, dnstype uint16,results chan string){
			defer wg.Done()
			resolve(ip, dnstype,results)
		} (query.ip, query.dnstype,results)

		//simulate the delay
		delay := r.NormFloat64() * StdDev + Mean
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	wg.Wait()
	close(results) //close channel
	fmt.Println("Writing results to file dat.csv........")

	//Write all the results into dat.cvs file
	f, _ := os.Create("/Users/Karen/workArea/charmander/experiments/nessy/services/normalload/dat.csv")
	for r := range results {
		f.WriteString(r)
	}
	f.Sync()
	f.Close()
	fmt.Println("All done! Try typing: cat dat.csv")
}

func resolve(ip string, dnstype uint16,results chan string ){
	message := new(dns.Msg)
	message.SetQuestion(ip,dnstype)
	client := new(dns.Client)
	client.DialTimeout = 7
	response, response_time, _ := client.Exchange(message, "172.31.2.12:53")
	s := "success"
	if response == nil {
		s = "failed"
	}

	fmt.Println(ip, response_time)
	str := ip + "\t" + s + "\t" + response_time.String() + "\n"
	results <- str //push into channel
}

func collectQueries(fileop string) []query {

	// Read file
	//f, err := os.Open("/data/"+fileop)
	f, err := os.Open("/Users/Karen/workArea/charmander/experiments/nessy/services/normalload/queryfiles/"+fileop)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err when opening file:", err)
	}

	//Create a slice to store these queries
	size := 100 //default
	if(fileop == "tiny_20"){size = 20} else if(fileop == "medium_1500") {size = 1500}
	queries := make([]query,size)

	//Scan file and split into words
	scanner := bufio.NewScanner(f)
	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	count := 0


	for scanner.Scan() {
		ip := scanner.Text()
		var dnstype uint16
		scanner.Scan()
		dnstype = type_to_uint(scanner.Text())
		queries[count] = query{ip, dnstype}
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	f.Close()

	return queries
}



func type_to_uint(str string) uint16{
	switch str{
	case  "A":
		return dns.TypeA
	case  "MX":
		return dns.TypeMX
	case  "PTR":
		return dns.TypePTR
	case "AAAA":
		return dns.TypeAAAA
	}
	return dns.TypeA
}
