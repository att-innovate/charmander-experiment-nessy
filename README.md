Charmander-Experiment: Nessy
----------------------------

#### Prerequisite
A local Charmander cluster has to be up and running.
Related documentation available at [https://github.com/att-innovate/charmander](https://github.com/att-innovate/charmander).

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


#### Start cAdvisor and Analytics-Stack

    ./bin/start_cadvisor
    ./bin/start_analytics

#### Start DNS server on slave2

    ./experiments/nessy/bin/start_dns-sl2


#### Start DNS performance tool (Nominum) on slave 3. 
We provide three sizes of queryfiles, the tiny one has 20 IP addresses, the small one has 100, and the medium has 1500. Type one of the following commands with your chosen size to send queries to the DNS server. 

	
	 ./experiments/nessy/bin/start_dnsperf-sl3-tiny

    
     ./experiments/nessy/bin/start_dnsperf-sl3-small


     ./experiments/nessy/bin/start_dnsperf-sl3-medium

#### Now let's check the DNS performance result
You can go to [Mesos](http://172.31.1.11:5050/#/), click on sandbox of the dnsperf-sl3-tiny(or small, or medium). In the sandbox, click stdout to see the output of dnsperf.

You can also check the CAdvisor for [DNS server (slave 2)](http://slave2:31500/containers/)  and [DNS performave tool (slave 3)](http://slave3:31500/containers/) to see their usage. 

#### That's it, let's clean up

    ./bin/reset_cluster

..and head back to the Charmander [Homepage](https://github.com/att-innovate/charmander/)
