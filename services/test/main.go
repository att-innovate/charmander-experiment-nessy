package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/miekg/dns"
	"strings"
	"sync/atomic"
	"flag"
)


var numUsers = flag.Int("u",10,"number of users(threads) sending queries, default is 10")
var varyNumUsers = flag.Int("vu",0,"number of users(threads) to increase or decrease by the end of varyTime, has to specify -vt varyTime and numUsers")
var varyTime = flag.Int("vt",0,"length of time that the number of users vary, unit is in minute")
var maxQueries = flag.Int("q",0,"max number of queries, default(or if you type 0) is infinite number of querys")
var tlimit = flag.Int("t",0,"max time limit, default(or if type 0) is infinite" )
var Mean = flag.Float64("m",2000.0,"mean of the client's query rate distribution" )
var StdDev = flag.Float64("d",1000.0,"standard deviation of the client's query rate distribution" )


var activeRoutines = new(int32)
var numQueries = new(int32)


func main() {
	flag.Parse()
	changePerMin := 0

	if *varyNumUsers!=0 && *varyTime==0 {
		fmt.Println("Must specify -t (the time limit) to use -v (vary number of users)")
		return
	}else if(*numUsers + *varyNumUsers < 0){
		fmt.Println("the number of users to decrease cannot exceed the original number of users")
		return
	}else if *varyTime!=0 {
		changePerMin = *varyNumUsers / *varyTime
		fmt.Println("the number of threads to vary per min (varyNumUsers/tlimit) is : ", changePerMin,"(0 means none)")
	}

	fmt.Println("-u: the number of users is : ", *numUsers)
	fmt.Println("-q: the number of max queries is : ", *maxQueries, "(0 means infinite)")
	fmt.Println("-t: the number of max time in seconds is : ", *tlimit,"(0 means infinite)")

	


	//Random Generater with seed 99
	r := rand.New(rand.NewSource(99))

	done := make(chan bool)
	timeup := make (chan bool)
	varyThread := make (chan bool)

	go func (){
		timer(timeup,varyThread)
	}()



	for pt := 0;  pt < *numUsers; pt++ {
		go func() {
			doIt(r,done,timeup)
		}()
		atomic.AddInt32(activeRoutines,1)
	}



	if *tlimit == 0 && *maxQueries == 0 && *varyNumUsers == 0{
		exit := make(chan bool) //wait forever!
		<-exit
	}else {
		for{
			select{
				case <- varyThread:
					
					for pt := 0;  pt < changePerMin; pt++ {
						go func() {
							doIt(r,done,timeup)
						}()
						atomic.AddInt32(activeRoutines,1)
						fmt.Println("Added a thread, now have activeRoutines: ", *activeRoutines)
					}


				case <- timeup :
					for *activeRoutines > 0 {
						done <- true
					}
					return

				default :
			}
		}
	}


}

func timer(timeup chan bool, varyThread chan bool){
	if *tlimit != 0 {
		for runtime:=0; runtime < *tlimit; runtime += 60 {
			time.Sleep(time.Duration(60) * time.Second)
			if(*varyNumUsers != 0) {
				varyThread <- true
			}
		}
	}else if *varyNumUsers != 0 {
		for {
			time.Sleep(time.Duration(60) * time.Second)
			varyThread <- true	

		}
	}else{
		return
	}
	
	fmt.Println(" Time up, after ", *tlimit, "seconds.")
	timeup <- true
}

func doIt ( r *rand.Rand, done chan bool, timeup chan bool) {

	client := new(dns.Client)
	client.DialTimeout = time.Duration(5) * time.Second
	client.ReadTimeout = time.Duration(5) * time.Second
	message := new(dns.Msg)


	lines:=[]string{"homer.ave.	MX",
					"homer.ave.	NS",
					"homer.ave.	A",
					"server.homer.ave.	A",
					"host1.homer.ave.	A",
					"host2.homer.ave.	CNAME",
					"host3.homer.ave.	A",
					"host3.homer.ave.	A",
					"mail.homer.ave.	CNAME",
					"server.homer.ave.	A",
					"www.homer.ave.	CNAME",

					"host2015.homer.ave.	AAAA", //random false one query
					}
					
	size := len(lines)



	for{

			select {
			case <- done:
				atomic.AddInt32(activeRoutines,-1)
				return

			default:
				delay := r.NormFloat64() * (*StdDev) + (*Mean)
				time.Sleep(time.Duration(delay) * time.Millisecond)

				tokens := strings.Split(lines[r.Intn(size)], "\t")

				message.SetQuestion(tokens[0], resolveDNSType(tokens[1]))

				_,_,_ = client.Exchange(message, "172.31.2.12:53")

				fmt.Println(tokens[0],tokens[1])

				atomic.AddInt32(numQueries,1)

				if *maxQueries!=0 && *numQueries >= int32(*maxQueries) {
					fmt.Println("Quitting. Have reached the max number of queries: ", *maxQueries)
					timeup <- true //notice the main thread and triggers the done for all the other threads
					atomic.AddInt32(activeRoutines,-1)
					return
				}


			}

	}


}

func resolveDNSType(str string) uint16{
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
}