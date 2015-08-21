Charmander-Experiment: Nessy
----------------------------

Nessy is an experiment that runs on the Charmander Lab Platform that performs the dynamic orchestration of DNS servers and detection of DDoS Attack. 


![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/Nessy.jpg?raw=true)


#### Main features of Nessy:
a) Allow for a dynamic orchestration of containerized DNS servers (configured with Bind9 services) 
    
b) Generating the necessary training data sets by implementing load generators for both normal user load & DDoS traffic patterns
    
c) Identifying & training & updating the appropriate machine learning model
    
d) (Not fully implemented yet) Processing data collected from DNS servers with learned model and providing intelligent suggestions to scheduler




1. [Set Up Nessy in Charmander Environment][1]

    - Configure and build the five nodes with Vagrant and VirtualBox
    - Reload and reset environment

2. [Build analyzers, vector and Nessy] [2] 
    
    - Build and configure the containers for nessy experiemnt
    - Build analytics, datacollectors and vector

Now, it is time to run and play with it! 

First, Open the [Mesos](http://172.31.1.11:5050/#/) to check tasks that are actively running.

Then, open up[Vector](http://http://172.31.2.11:31790/#/?host=slave2&hostspec=localhost) to see the metrics collected from DNS server on slave2.

#### Run a script we provide to check what different load patterns will affect the DNS server

These scripts will automatically set up everything and run some normalloads and DDos attacks to the DNS server. 
Note:it does not include smartscaling.

    ./experiments/nessy/bin/script_2min
    ./experiments/nessy/bin/script_120min
    .....


#### Alternatively, you can run every step manuelly. Follow the steps below:

#### Start Vector and Analytics-Stack
  
    ./bin/start_analytics
    ./bin/start_vector

#### Start DNS server on slave2

    ./experiments/nessy/bin/start_dns-sl2


#### Start smartscaling 
Smartscaling is the machine intelligence pieve that implements autoscaling and DDos detection. After set up, smartscaling will tell Charmander scheduler to spin up /shut down another DNS server on demand. It can also kill the DDoS attack once detected.

	
	 ./experiments/nessy/bin/start_smartscaling


#### Start Normal Load Emulation with various loads
This will start a load emulation where 200 users simutanously send queries to the DNS server in a random pattern

    ./experiments/nessy/bin/start_normalload-u200


After a while, you can put on a heavier load on the server. Say, 300 more users.
    
    ./experiments/nessy/bin/start_normalload-u100
    ./experiments/nessy/bin/start_normalload-u100
    ./experiments/nessy/bin/start_normalload-u100

At this point, you should see in Vector that the network throughput of the DNS server goes up quickly when serving 500 users. Soon, you would observe in Mesos that another DNS server in slave 3 (dns-sl3) appears. It is scaled up by scheduler according to the smartscaling, and it will take over about half of the load from the previous server. You can also check the metrics for this new DNS server on [Vector](http://172.31.2.11:31790/#/?host=slave3&hostspec=localhost)


#### You can now kill the 300 users.

    ./bin/kill_task normalload-u100

Soon, you will observe that dns-sl3 disappears, being shut down by schedule as smartscaling realize that we don't need two servers runnning to serve only 200 users. Now, you should observe in Vector that network throughput of dns-sl2 will go down first (#due to lowered users traffic#) and then goes up a little bit (#take over all the load from dns-sl3#).

#### While Normal Load is running, you can also start a DDos attack.

This will start a Bonesi DDos attack to the DNS server. ([Bonesi](https://github.com/markus-go/bonesi) is an open-source botnet simulator developed by Markus Goldstein.)

    ./experiments/nessy/bin/start_bonesi-500
    
This will put a lot of load on the network of DNS server. However, smartscaling should be able to detect it is an attack rather than a scaling-up situation. THen you should be observing the bonesi-500 being killed and the network throughput returns normal.

#### Now let's check the data collected from InfluxDB

// to be implemented

 

#### That's it, let's clean up

    ./bin/reset_cluster

..and head back to the Charmander [Homepage](https://github.com/att-innovate/charmander/)


[1]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/SETUPNESSYNODES.md
[2]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/BUILDNESSY.md
