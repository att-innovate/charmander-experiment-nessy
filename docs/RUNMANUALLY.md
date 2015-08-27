Run Tasks manually 
----------------------

#### Start Vector and Analytics-Stack
  
    ./bin/start_analytics
    ./bin/start_vector

#### Start DNS server on slave2

    ./experiments/nessy/bin/start_dns-sl2


#### Start smartscaling 
Smartscaling is the machine intelligence pieve that implements autoscaling and DDos detection. After set up, smartscaling will tell Charmander scheduler to spin up /shut down another DNS server on demand. It can also kill the DDoS attack once detected.
	
	 ./experiments/nessy/bin/start_smartscaling

Then, open the [Mesos][2] to check tasks that are actively running. It should look like this:

![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/MesosExp.png?raw=true)


To see the data generated on run-time, open the [Vector][3] to observe the metrics collected from DNS server on slave2.


#### Increase traffic 
This will start a load emulation where 200 users simutanously send queries to the DNS server in a random pattern, for example, normal load of 200 users

    ./experiments/nessy/bin/start_normalload-u200


After a while, you can put on a heavier load on the server. Say, 300 more users.
    
    ./experiments/nessy/bin/start_normalload-u100
    ./experiments/nessy/bin/start_normalload-u100
    ./experiments/nessy/bin/start_normalload-u100

At this point, you should see in Vector that the network throughput of the DNS server goes up quickly when serving 500 users. 

Soon, you would observe in Mesos that another DNS server in slave 3 (dns-sl3) appears. It is scaled up by scheduler according to the smartscaling, and it will take over about half of the load from the previous server. 
You can also check the metrics for this new DNS server on [Vector][3]


#### Lowering the traffic

    ./bin/kill_task normalload-u100

Essentialy, this command kills the simulation of 300 users.

Soon, you will observe in Mesos that dns-sl3 disappears, because it is shut down by scheduler as smartscaling realize that we don't need two servers runnning. Meantime, you should observe in Vector that network throughput of dns-sl2 will go down first (#due to lowered users traffic#) and then goes up a little bit (#take over all the load from dns-sl3#).

#### Start a DDos attack.

You can start a DDoS attack while the normal load is running. This will start a Bonesi DDos attack to the DNS server. ([Bonesi](https://github.com/markus-go/bonesi) is an open-source botnet simulator developed by Markus Goldstein.)

    ./experiments/nessy/bin/start_bonesi-500
    
This will put a lot of load on the network of DNS server, which might first appears to be similar to a heavy load. However, smartscaling should be able to detect it is an attack rather than a scaling-up situation. 

You will see in [Mesos][2] that the bonesi-500 being killed and in Vector that the network throughput returns normal.


####Collect and Explore the data
When the script complets, you can collect all the data it generated from InfluxDB into a csv file by querying from Spark. (InfluxDB itself does not provide csv conversion)
In your SPARK_FOLDER, run:

	./bin/spark-shell --jars $(echo ./lib/*.jar | tr ' ' ',')  

Then inside the spark-shell (Scala environment), run:

	>>> import org.att.charmander.CharmanderUtils
	>>> CharmanderUtils.exportDataToCSV(CharmanderUtils.VECTOR_DB,[YOUR_QUERY],[YOUR_OUTPUT_PATH])

Examples for [YOUR_QUERY]: "select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave2' limit 10000"). For more query language, please refer to [InfluxDB Web UI][1]


If you wish to play with the data inside the Spark shell, we provide functions to turn the data into RDD (Resilient Distributed Dataset) for you to use. Please refer to [charmender-spark ChamanderUtils][4]]
Note: after resetting the cluster, you won't be able to retrieve the data from previous experiments if that are not saved in an output file.

If you wish, you can further implement the codes you plays with in spark-shell into the analytics piece, please refer to [You can build your own Nessy analytics][5]


#### When you are done with experiment, clean up the cluster

    ./bin/reset_cluster


[1]: https://influxdb.com/docs/v0.8/introduction/getting_started.html
[2]: http://172.31.1.11:5050/#/
[3]: http://172.31.2.11:31790/#/?host=slave3&hostspec=localhost
[4]: https://github.com/att-innovate/charmander-spark/blob/master/src/main/scala/org/att/charmander/CharmanderUtils.scala
[5]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/analytics/README