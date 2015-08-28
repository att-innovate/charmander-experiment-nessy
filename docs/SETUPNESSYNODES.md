Set Up Nessy in Charmander Environment
--------------------------------------




![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/Nessy_Implementation.jpg?raw=true)

#### Install Charmander and Nessy
First of all, you will need to install Charmander. 
In your working directory:

    git clone https://github.com/att-innovate/charmander.git

Next, install Nessy inside the experiments folder under Charmander
    
    cd charmander/experiments
    git clone https://github.com/att-innovate/charmander-experiment-nessy.git nessy


#### Set Up the nodes in Charmander to run Nessy

Nessy requires a different cluster configuration from Charmander's default setting. 
Overwrite the cluster.yml in charmander with the cluster.yml file in nessy folder:
    
    cd .. 
    cp ./experiments/nessy/cluster.yml ./


Check whether the cluster.yml is correct
	
	cat ./cluster.yml

It should look exactly like this:

	# Mesos cluster configurations

	# The numbers of slaves,
	# includes node for analytics and optional traffic generator
	# hostname will be slave1,slave2,…
	############################################################
	slave_n : 4

	# Name of slave-node assigned to do analytics
	# All other slave nodes become "lab-nodes"
	################################################
	analytics_node : "slave1"

	# Name of slave-node assigned to do traffic generators
	# If left empty no traffic node gets set
	################################################
	traffic_node : "slave4"

	# Memory and Cpu settings
	##########################################
	master_mem     : 384
	master_cpus    : 1
	analytics_mem  : 2048
	analytics_cpus : 3
	slave_mem      : 1024
	slave_cpus     : 2
	traffic_mem    : 1024
	traffic_cpus   : 2

	# Private ip bases, don't change
	################################################
	master_ipbase  : "172.31.1."
	slave_ipbase   : "172.31.2."


Now, we can set up the nodes for Nessy. (If you have charmander set up before, you will need to do vagrant destroy first and then rebuild it)

Please refer to the Charmander's [Setup Nodes](https://github.com/att-innovate/charmander/blob/master/docs/SETUPNODES.md) and follow the instructions on that page.





Verify that Mesos is connected and no task is running using the Mesos console at [http://172.31.1.11:5050](http://172.31.1.11:5050)



Cool, we have the Nessy nodes all set up!


Next step: [Build analyzers, vector and Nessy] [2] 


[2]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/BUILDNESSY.md