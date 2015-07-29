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

var numUsers = flag.Int("u",20,"number of users(threads) sending queries, default is 10")
var maxQueries = flag.Int("q",0,"max number of queries, default(or if you type 0) is infinite number of querys")
var tlimit = flag.Int("t",0,"max time limit, default(or if type 0) is infinite" )


var numFailed = new(int32)
var numResponse = new(int32)

const StdDev = 300.0
const Mean = 500.0


type Resolve struct {
	ip string
	dnstype uint16
	WP *workpool.WorkPool
	rg *rand.Rand
}


func main() {
	flag.Parse()
	fmt.Println("the number of users is : ",*numUsers)
	fmt.Println("the number of max queries is : ", *maxQueries, "(0 means infinite)")
	fmt.Println("the number of max time in seconds is : ", *tlimit,"(0 means infinite)")

	//read the file and scan them into words
	queryfile, err := os.Open("/normalload/queryfiles/medium_1500")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err when opening file:", err)
	}


	//Random Generater
	r := rand.New(rand.NewSource(99))

	workPool := workpool.New(*numUsers,10000)

	shutdown := false

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

				_ = workPool.PostWork("routine", &work)



				count++
				//quit of it reaches max number of queries or if it is shut down
				if (*maxQueries!=0 && count >= *maxQueries) || shutdown == true {
					fmt.Println("totol count is: ",count)
					return
				}

				//if the queue is too long, wait for one sec before assign queires to workers
				if(workPool.QueuedWork() >= 8000){
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	if *maxQueries!=0 {
		fmt.Println("Number of queries has reached to maximum: ", *maxQueries)
	}else if *tlimit == 0 {
		//fmt.Println("Running time is set to be infinite, you can hit ENTER to exit.....")
		//reader := bufio.NewReader(os.Stdin)
		//reader.ReadString('\n')
		exit := make(chan bool) //wait forever!
		<-exit
	}else{
		//sleep til time up
		time.Sleep(time.Duration(*tlimit) * time.Second)
		fmt.Println(" Time up, after ", *tlimit, "seconds.")
	}


	shutdown = true

	fmt.Println("Shutting Down")
	workPool.Shutdown("routine")

	queryfile.Close()
	fmt.Println("All done! Running time: ",*tlimit,"There are ", *numFailed ," failures in ", *numResponse," responses. ")

}

func (rs *Resolve) DoWork(workRoutine int) {
	//simulate the delay, with normal distribution
	delay := rs.rg.NormFloat64() * StdDev + Mean
	time.Sleep(time.Duration(delay) * time.Millisecond)
	message := new(dns.Msg)
	message.SetQuestion(rs.ip,rs.dnstype)
	client := new(dns.Client)
	client.DialTimeout = 10000000
	response ,response_time, _ := client.Exchange(message, "172.31.2.12:53")
	atomic.AddInt32(numResponse,1)
	if response == nil {
		atomic.AddInt32(numFailed,1)
	}

	fmt.Println("[",workRoutine,"]",rs.ip, response_time)

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
