package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/miekg/dns"
	//"strings"
	"sync/atomic"
	"flag"
)

type query struct{
	ip string
	dnstype uint16
}


//The 2 DNS servers address. When they are both up, queries will be equally distributed to them
type dnsserver struct {
    addr string
    up  bool
}


var dns1 = dnsserver{"172.31.2.12:53",false}
var dns2 = dnsserver{"172.31.2.14:53",false}
var servers = [...]dnsserver{dns1,dns2} //array of servers, in this experiment we have 2 dns servers


//Advanced settings for porportionally increasing/decreasing users during a period of time. 
var varyNumUsers = flag.Int("vu",0,"You can specify the number of users(threads) to increase or decrease by the end of varyTime, has to specify -vt varyTime and numUsers")
var varyTime = flag.Int("vt",0,"You can specify the period of time (the number of varyIntervals) in which the user number changes ")
var varyInterval = flag.Int("vi", 60, "vary the number of users after every interval, default is every minute")

var numUsers = flag.Int("u",10,"number of users(threads) sending queries, default is 10")
var maxQueries = flag.Int("q",0,"max number of queries, default(or if you type 0) is infinite number of querys")
var tlimit = flag.Int("t",0,"max time limit, default(or if type 0) is infinite" )
var Mean = flag.Float64("m",2000.0,"mean of the client's query rate distribution" )
var StdDev = flag.Float64("d",1000.0,"standard deviation of the client's query rate distribution" )
var verbose = flag.Bool("verbose",false,"Print out information about each query and the number of threads" )



var activeRoutines = new(int32)
var numQueries = new(int32)



func main() {
	flag.Parse()
	changePerInterval := 0

	if *varyNumUsers!=0 && *varyTime==0 {
		if *verbose {
			fmt.Println("Must specify -vt (the time limit) to use -vu (vary number of users)")
		}
		return
	}else if(*numUsers + *varyNumUsers < 0){
		if *verbose {
			fmt.Println("number of users is", *numUsers, ", cannot decrease users to a negative number.")
		}
		return
	}else if *varyTime!=0 {
		changePerInterval = *varyNumUsers / *varyTime
		if *verbose {
			fmt.Println("the number of threads to vary per varyInterval (varyNumUsers / varyTime) will be:", changePerInterval)
			fmt.Println("The time period of increasing/decreasing users (varyInterval * varyTime) will be:", *varyInterval * (*varyTime))
		}
	}


	//Random Generater with seed 99
	r := rand.New(rand.NewSource(99))

	done := make(chan bool)
	timeup := make (chan bool)
	varyThread := make (chan bool)

	go func (){
		checkServers(done)
	}()
	
	go func (){
		timer(timeup,varyThread)
	}()
	

	for pt := 0;  pt < *numUsers; pt++ {
		go func() {
			doIt(r,done,timeup)
		}()
	}



	if *tlimit == 0 && *maxQueries == 0 && *varyNumUsers == 0{
		exit := make(chan bool) //wait forever!
		<-exit
	}else {
		for{
			select{
				case <- varyThread:
					if changePerInterval > 0 {
						for pt := 0;  pt < changePerInterval; pt++ {
							go func() {
								doIt(r,done,timeup)
							}()
						}
					}else if changePerInterval < 0{
						for pt := 0;  pt > changePerInterval; pt-- {
							done <- true
						}
					}
					
				case <- timeup :
					for *activeRoutines > 0 { //also need to terminate checkServer
						done <- true
					}
					return

				default :
			}
		}
	}


}

func checkServers(done chan bool){
	client := new(dns.Client)
	client.DialTimeout = time.Duration(5) * time.Second
	client.ReadTimeout = time.Duration(5) * time.Second
	message := new(dns.Msg)

	for{

			message.SetQuestion("homer.ave.", dns.TypeA)	
			for i := range servers{
				res,_,_ := client.Exchange(message, servers[i].addr)	
				servers[i].up = (res!=nil)
				
				if *verbose {
					fmt.Println("dns-slv",i+2,"up is: ",servers[i].up)
				}
			}		
			
			time.Sleep(time.Duration(6) * time.Second) //check every minute
		
	}

}


func timer(timeup chan bool, varyThread chan bool){
	runtime:=0

	if *varyNumUsers != 0 {
		for runtime < *varyTime * (*varyInterval) {
			time.Sleep(time.Duration(*varyInterval) * time.Second)
			runtime += *varyInterval 
			varyThread <- true

			if *tlimit!=0 && runtime > *tlimit - *varyInterval {
				time.Sleep(time.Duration(*tlimit - runtime) * time.Second)
				fmt.Println(" Time up, after ", *tlimit, "seconds.")
				timeup <- true
			}
		}

		if *verbose {
			fmt.Println("End varying users.")
		}
	} 

	if *tlimit != 0 {
		time.Sleep(time.Duration(*tlimit - runtime) * time.Second)
		fmt.Println(" Time up, after ", *tlimit, "seconds.")
		timeup <- true
	}

}

func doIt ( r *rand.Rand, done chan bool, timeup chan bool) {
	atomic.AddInt32(activeRoutines,1)
	if *verbose {
		fmt.Println("Added a thread, now have activeRoutines: ", *activeRoutines)
	}
	
	client := new(dns.Client)
	client.DialTimeout = time.Duration(5) * time.Second
	client.ReadTimeout = time.Duration(5) * time.Second
	message := new(dns.Msg)

	queries:=[]query{ {"homer.ave.", dns.TypeMX},
					{"homer.ave.", dns.TypeNS},
					{"homer.ave.", dns.TypeA},
					{"server.homer.ave.", dns.TypeA},
					{"host1.homer.ave.", dns.TypeA},
					{"host2.homer.ave.", dns.TypeCNAME},
					{"host3.homer.ave.", dns.TypeA},
					{"host3.homer.ave.", dns.TypeA},
					{"mail.homer.ave.", dns.TypeCNAME},
					{"server.homer.ave.", dns.TypeA},
					{"www.homer.ave.", dns.TypeCNAME},
					{"host2015.homer.ave.", dns.TypeAAAA}, //random false one query
					}
					
	size := len(queries)
	//addr := nil

	for{
		select {
		case <- done:
			if *verbose {
				fmt.Println("Deleteed a thread, now have activeRoutines: ", *activeRoutines)
			}
			atomic.AddInt32(activeRoutines,-1)
			return

		default:
			for i := range servers{
				if !servers[i].up {
					continue      // If this server is down, go check next server
				}

				delay := r.NormFloat64() * (*StdDev) + (*Mean)
				time.Sleep(time.Duration(delay) * time.Millisecond)
				
				//tokens := strings.Split(lines[r.Intn(size)], "\t")
				randomQuery := queries[r.Intn(size)]
				message.SetQuestion(randomQuery.ip, randomQuery.dnstype)
				
			
				_,response_time,_ := client.Exchange(message, servers[i].addr)			
				
				if *verbose{
					fmt.Println("server",i," ",randomQuery.ip,response_time)
				}

				atomic.AddInt32(numQueries,1)
				if *maxQueries!=0 && *numQueries >= int32(*maxQueries) {
					fmt.Println("Quitting. Have reached the max number of queries: ", *maxQueries)
					timeup <- true //notice the main thread and triggers the done for all the other threads
					if *verbose {
						fmt.Println("Deleteed a thread, now have activeRoutines: ", *activeRoutines)
					}
					atomic.AddInt32(activeRoutines,-1)
					return
				}
			}
			

		}

	}


}

/*func resolveDNSType(str string) uint16{
	switch str{
	case  "A":
		return dns.TypeA
	case  "MX":
		return dns.TypeMX
	case  "PTR":
		return dns.TypePTR
	case "AAAA":
		return dns.TypeAAAA
    case "CNAME":
		return dns.TypeCNAME
}
	return dns.TypeA
}*/


