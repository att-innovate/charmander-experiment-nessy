Charmander-Experiment: Nessy
----------------------------

Nessy is an experiment that runs on the Charmander Lab Platform that performs the dynamic orchestration of DNS servers and detection of DDoS Attack. 

![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/Nessy_Implementation.jpg?raw=true)

#### Main features of Nessy:
    a) Allow for a dynamic orchestration of containerized DNS servers (configured with Bind9 services) 
    
    b) Generating the necessary training data sets by implementing load generators for both normal user load & DDoS traffic patterns
    
    c) Identifying & training & updating the appropriate machine learning model
    
    d) (Not fully implemented yet) Processing data collected from DNS servers with learned model and providing intelligent suggestions to scheduler


#### Prerequisite

First, a Charmander cluster needs to up and running before we install Nessy. In your working directory:

    git clone https://github.com/att-innovate/charmander.git
    cd charmander/experiments
    git clone https://github.com/att-innovate/charmander-experiment-nessy.git

Since Nessy requies 4 nodes rather than the default 3 nodes in Charmander, you will need to overwrite the cluster.yml in charmander with the cluster.yml for Nessy:
    
    cd .. 
    cp /experiments/nessy/cluster.yml ./

Now, we can set up the nodes in Charmander clusters. Please refer to the related documentation available at [https://github.com/att-innovate/charmander](https://github.com/att-innovate/charmander).


Verify that you are in your local Charmander directory and reset the Charmander cluster.

    ./bin/reset_cluster

Verify that no task is running using the Mesos console at [http://172.31.1.11:5050](http://172.31.1.11:5050)

#### Build and deploy the simulators and the analyzer

Lets build it first. Change to the experiments folder and check out the code into a folder called `nessy`

    cd experiments
    git clone https://github.com/att-innovate/charmander-experiment-nessy.git nessy

Change your working directory back to the root of Charmander and start the build process

    cd ..
    ./experiments/nessy/bin/build

This command builds and creates and deploys Docker images for nessy.
This process will take some time the first time you run it.

#### You can run a nessy_script which will automatically set up everything and  includes some normalloads and one DDos attack to the DNS server

    ./experiments/nessy/bin/script_2min


Alternatively, you can run every step manuelly. Follow the steps below:

#### Start Vector and Analytics-Stack

    ./bin/start_vector
    ./bin/start_analytics

#### Start DNS server on slave2

    ./experiments/nessy/bin/start_dns-sl2


#### Start DNS performance tool (Nominum) on slave 3. 
We provide three sizes of queryfiles, the tiny one has 20 IP addresses, the small one has 100, and the medium has 1500. Type one of the following commands with your chosen size to send queries to the DNS server. 

	
	 ./experiments/nessy/bin/start_dnsperf-sl3-tiny

    
     ./experiments/nessy/bin/start_dnsperf-sl3-small


     ./experiments/nessy/bin/start_dnsperf-sl3-medium

#### Start Normal Load Emulation
This will start a load emulation where 100 users send queries to the DNS server in a random pattern

    ./experiments/nessy/bin/start_normalload-u100


#### While Normal Load is running, you can also start a DDos attack
This will start a Bonesi DDos attack to the DNS server

    ./experiments/nessy/bin/start_bonesi-500
    
To stop the DDos attack, you can run

    ./bin/kill_task bonesi

#### Now let's check the data from DNS server
You can go to [Mesos](http://172.31.1.11:5050/#/), click on sandbox of the dnsperf-sl3-*, or normalload-* ). In the sandbox, click stdout to see the output of dnsperf.

You can also check the Vector for [DNS server (slave 2)](http://172.31.2.11:31790/#/?hostspec=localhost&host=172.31.2.12)  and [DNS performave tool (slave 3)](http://slave3:31500/containers/) to see their usage. 
 

#### That's it, let's clean up

    ./bin/reset_cluster

..and head back to the Charmander [Homepage](https://github.com/att-innovate/charmander/)
