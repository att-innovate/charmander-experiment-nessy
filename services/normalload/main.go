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
var maxQueries = flag.Int("q",0,"max number of queries, default(or if you type 0) is infinite number of querys")
var tlimit = flag.Int("t",0,"max time limit, default(or if type 0) is infinite" )
var Mean = flag.Float64("m",2000.0,"mean of the client's query rate distribution" )
var StdDev = flag.Float64("d",1000.0,"standard deviation of the client's query rate distribution" )


var activeRoutines = new(int32)
var numQueries = new(int32)


func main() {
	flag.Parse()
	fmt.Println("-u: the number of users is : ", *numUsers)
	fmt.Println("-q: the number of max queries is : ", *maxQueries, "(0 means infinite)")
	fmt.Println("-t: the number of max time in seconds is : ", *tlimit,"(0 means infinite)")

	//Random Generater with seed 99
	r := rand.New(rand.NewSource(99))

	done := make(chan bool)
	timeup := make (chan bool)

	if *tlimit!= 0 {
		go func (){
			timer(timeup)
		}()
	}



	for pt := 0;  pt < *numUsers; pt++ {
		go func() {
			doIt(r,done,timeup)
		}()
	}



	if *tlimit == 0 && *maxQueries == 0 {
		exit := make(chan bool) //wait forever!
		<-exit
	}else {
		for{
			select{
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

func timer(timeup chan bool){
	time.Sleep(time.Duration(*tlimit) * time.Second)
	fmt.Println(" Time up, after ", *tlimit, "seconds.")
	timeup <- true
}

func doIt ( r *rand.Rand, done chan bool, timeup chan bool) {

	client := new(dns.Client)
	client.DialTimeout = time.Duration(5) * time.Second
	client.ReadTimeout = time.Duration(5) * time.Second
	message := new(dns.Msg)


	lines:=[]string{"homer.ave.	MX",
					"server.homer.ave.	A",
					"host1.homer.ave.	A",
					"host2.homer.ave.	A",
					"host3.homer.ave.	A",
					"mail.homer.ave.	A",
					"server.homer.ave.	A",
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


