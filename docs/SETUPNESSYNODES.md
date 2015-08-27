Set Up Nessy in Charmander Environment
--------------------------------------

![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/Nessy_Implementation.jpg?raw=true)

#### Step 1 - Install Charmander
If you have already have a Charmander set up, the 4-node charmander cluser will need to be reconstructed into a 5-node cluster to run Nessy. You will need to do:

	vagrant -y destroy

Then skip the following and go directlly to Step 2.

Otherwise, a Charmander cluster needs to be first up and running before we install Nessy. In your working directory:

    git clone https://github.com/att-innovate/charmander.git

#### Step 2 - Install Nessy

Nest, we will need to install nNessy inside the experiments folder under Charmander
    
    cd charmander/experiments
    git clone https://github.com/att-innovate/charmander-experiment-nessy.git nessy


#### Step 3 - Set Up the nodes in Charmander to run Nessy
Since Nessy requies 5 nodes rather than the default 4 nodes in Charmander, you will need to overwrite the cluster.yml in charmander with the cluster.yml for Nessy:
    
    cd .. 
    cp /experiments/nessy/cluster.yml ./

Now, we can set up the nodes in Charmander clusters. Please refer to the related documentation available at [https://github.com/att-innovate/charmander](https://github.com/att-innovate/charmander).


Verify that you are in your local Charmander directory and reset the Charmander cluster.

    ./bin/reset_cluster

Verify that no task is running using the Mesos console at [http://172.31.1.11:5050](http://172.31.1.11:5050)