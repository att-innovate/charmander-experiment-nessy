Set Up Nessy in Charmander Environment
--------------------------------------

![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/Nessy_Implementation.jpg?raw=true)

#### Install Charmander and Nessy
You will need to install a Charmander cluster first. In your working directory:

    git clone https://github.com/att-innovate/charmander.git

Nest, we will need to install nNessy inside the experiments folder under Charmander
    
    cd charmander/experiments
    git clone https://github.com/att-innovate/charmander-experiment-nessy.git nessy


#### Set Up the nodes in Charmander to run Nessy
Since Nessy requies 5 nodes rather than the default 4 nodes in Charmander, you will need to overwrite the cluster.yml in charmander with the cluster.yml for Nessy:
    
    cd .. 
    cp /experiments/nessy/cluster.yml ./

Now, we can set up the nodes for Nessy. (If you have charmander set up before, you will need to destroy it first.)
Please refer to the related documentation available at [Setup Nodes](https://github.com/att-innovate/charmander).


Verify that you are in your local Charmander directory and reset the Charmander cluster.

    ./bin/reset_cluster

Verify that no task is running using the Mesos console at [http://172.31.1.11:5050](http://172.31.1.11:5050)


Next step: [Build analyzers, vector and Nessy] [2] 


[2]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/BUILDNESSY.md