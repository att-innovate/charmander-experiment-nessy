Run A Script
------------

#### What's inside 

We have provided some scripts to minimize your work. It outs various patterns of load on one DNS server. They are handy to collect data for analysis or building models. 

In a script:
  
  - The script will automatically start vector, analytics and a DNS server 
  - You need to confirm that these tasks are up and running in Mesis.
  - Then it will run a series of normalloads and DDos attacks to the DNS server. It does NOT include smartscaling, and there will be only one DNS server.
  - The script will output a logfile.csv in your local charmander folder as a reference. It includes the number of users and whether there is a DDoS attacks(bonesi). It also records the timestamps showing the starting time of a each simulation that causes the change.it will appear in your local charmadner folder. 
  - You will need to manually collect and explore data generated during the experiment, as shown in the following detailed instructions.

#### Run the script

In your Charmander folder:

    ./experiments/nessy/bin/script_labels_30min


Open the [Mesos][2] to check tasks are all actively running. It should have all the tasks listed (order does not matter) shown in this picture:

![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/MesosExp.png?raw=true)


If you wish to see the data generated on run-time, open the [Vector][3] to observe the metrics collected from DNS server on slave2.

#### Collect and Explore the data

When the script complets, you can collect all the data it generated from InfluxDB into a csv file by querying from Spark. ([Instructions for instaling Apache Spark](http://spark.apache.org/downloads.html))

You will need  CharmanderUtils libraries. (For details, refer to [charmender-spark ChamanderUtils][4])

In your SPARK_FOLDER, run:
	
	cp /YOUR_WORKING_FOLDER/charmander/experiments/nessy/analytics/smartscaling/lib/charmander-utils_2.10-2.0.jar ./lib/

	./bin/spark-shell --jars $(echo ./lib/*.jar | tr ' ' ',')   //start a spark-shell that supports Scala 

Then inside the spark-shell (Scala environment), run:

	>>> import org.att.charmander.CharmanderUtils
	>>> CharmanderUtils.exportDataToCSV(CharmanderUtils.VECTOR_DB,[YOUR_QUERY],[YOUR_OUTPUT_PATH])

Examples for [YOUR_QUERY]: "select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave2' limit 10000"). For more query languages, please refer to [InfluxDB Web UI][1]


If you wish to play with the data inside the Spark shell, we provide functions to turn the data into RDD (Resilient Distributed Dataset) for you to use. Please refer to [charmender-spark ChamanderUtils][4]

Note: after resetting the cluster, you won't be able to retrieve the data from previous experiments if they are not saved in a file.


#### What's More

While the script is handy for running experiments and collecting data for training, we also provide a more flexible way to play with Nessy, and there can be more features like auto-scaling and DDoS deteciton. Please refer to [Manually starting tasks] [6]


If you wish, you can further implement the codes you plays with in spark-shell into the analytics piece, please refer to [You can build your own Nessy analytics][5]


[1]: https://influxdb.com/docs/v0.8/introduction/getting_started.html
[2]: http://172.31.1.11:5050/#/
[3]: http://172.31.2.11:31790/#/?host=slave3&hostspec=localhost
[4]: https://github.com/att-innovate/charmander-spark/blob/master/src/main/scala/org/att/charmander/CharmanderUtils.scala
[5]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/analytics/
[6]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/RUNMANUALLY.md

	