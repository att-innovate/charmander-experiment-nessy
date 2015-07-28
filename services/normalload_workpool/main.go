package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/miekg/dns"
	"os"
	"bufio"

	"github.com/goinggo/workpool"
    "sync/atomic"
	"flag"

)




//var fileop = flag.String("d", "tiny_20", "datafile, default is tiny_20")

var numUsers = flag.Int("u",10,"number of users(threads) sending queries, default is 10")
var maxQueries = flag.Int("q",0,"max number of queries, default(or if you type 0) is infinite number of querys")
var tlimit = flag.Int("-l",0,"max time limit, default(or if type 0) is infinite" )


var numFailed = new(int32)
var numQueries = new(int32)

const StdDev = 300.0
const Mean = 500.0


type Resolve struct {
	ip string
	dnstype uint16
	WP *workpool.WorkPool
	rg *rand.Rand
}


func (rs *Resolve) DoWork(workRoutine int) {
	message := new(dns.Msg)
	message.SetQuestion(rs.ip,rs.dnstype)
	client := new(dns.Client)
	client.DialTimeout = 10000000
	response , response_time, _ := client.Exchange(message, "172.31.2.12:53")
	atomic.AddInt32(numQueries,1)
	if response == nil {
		atomic.AddInt32(numFailed,1)
	}

	fmt.Println(rs.ip, response_time)


	//simulate the delay, with normal distribution
	delay := rs.rg.NormFloat64() * StdDev + Mean
	time.Sleep(time.Duration(delay) * time.Millisecond)



	//fmt.Println("QueuedWork: ",  mw.WP.QueuedWork(),  "ActiveRoutines: ", mw.WP.ActiveRoutines())

}



func main() {
	flag.Parse()
	fmt.Println("the number of users is : ",*numUsers)
	fmt.Println("the number of max queries is : ", *maxQueries, "(0 means infinite)")
	fmt.Println("the number of max time in seconds is : ", *tlimit,"(0 means infinite)")

	//read the file and scan them into words
	queryfile, err := os.Open("/Users/Karen/workArea/charmander/experiments/nessy/services/dnsperf/queryfiles/medium_1500")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err when opening file:", err)
	}



	r := rand.New(rand.NewSource(99))


	workPool := workpool.New(*numUsers,10000)

	shutdown := false

	timeout := make(chan bool, 1)
	/*if *tlimit != 0 {
		go func(){

			time.Sleep( time.Duration(*tlimit) * time.Second)
			timeout <- true
		}()
	}*/



	go func() {
		count := 0
		for {
			_,_ = queryfile.Seek(0,0)
			scanner := bufio.NewScanner(queryfile)
			// Set the split function for the scanning operation.
			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				ip := scanner.Text()
				scanner.Scan()

				dnstype := type_to_uint(scanner.Text())


				work := Resolve{
					ip: ip,
					dnstype: dnstype,
					WP: workPool,
					rg : r,
				}

				if err := workPool.PostWork("routine", &work); err != nil {
					fmt.Println("ERROR: %s\n", err)
					time.Sleep(100 * time.Millisecond)
				}

				count++

				if (*maxQueries!=0 && count >= *maxQueries) || shutdown == true {
					fmt.Println("totol count is: ",count)
					return
				}
			}
			for{
				if(workPool.QueuedWork() >= 5000){
					time.Sleep(100 * time.Millisecond)
				}else{
					break
				}
			}

		}




	}()

	fmt.Println("Hit ENTER to exit.....")

	select {
		case  <- timeout:
			fmt.Println("Time out, after ", *tlimit, "seconds")
		default:
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
	}



	shutdown = true

	fmt.Println("Shutting Down")
	workPool.Shutdown("routine")


	queryfile.Close()
	fmt.Println("All done! Time lasts: xxx There are ", *numFailed ," failures in ", *numQueries," queries. ")

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

/*
func collectQueries() []query {

	// Read file
	f, err := os.Open("/Users/Karen/workArea/charmander/experiments/nessy/services/dnsperf/queryfiles/small_100")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err when opening file:", err)
	}

	//Create a slice to store these queries
	queries := make([]query,100)

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
 */