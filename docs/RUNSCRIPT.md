Run a Nessy Script
------------------

#### Inside the script
These scripts will automatically set up everything and run a series of normalloads and DDos attacks to the DNS server. It does not include smartscaling.

The script will output a logfile.csv. It includes the number of users and whether there is a DDoS attacks(bonesi), and the timestamps showing the starting time of a simulation that causes the change.it will appear in your local charmadner folder. They can be a reference if you wish to run analysis based on these information.

#### Run the script

In your Charmander folder:

    ./experiments/nessy/bin/script_labels_30min

You will need to wait for some time to ensure all containers are up and running, and the data_collector is working.

Then, open the [Mesos][2] to check tasks that are actively running. It should look like this:

![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/MesosExp.png?raw=true)


If you wish to see the data generated on run-time, open the [Vector][3] to observe the metrics collected from DNS server on slave2.

####Collect and Explore the data
When the script complets, you can collect all the data it generated from InfluxDB into a csv file by querying from Spark. (InfluxDB itself does not provide csv conversion)
In your SPARK_FOLDER, run:

	./bin/spark-shell --jars $(echo ./lib/*.jar | tr ' ' ',')  

Then inside the spark-shell (Scala environment), run:

	>>> import org.att.charmander.CharmanderUtils
	>>> CharmanderUtils.exportDataToCSV(CharmanderUtils.VECTOR_DB,[YOUR_QUERY],[YOUR_OUTPUT_PATH])

Examples for [YOUR_QUERY]: "select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave2' limit 10000"). For more query language, please refer to [InfluxDB Web UI][1]


If you wish to play with the data inside the Spark shell, we provide functions to turn the data into RDD (Resilient Distributed Dataset) for you to use. Please refer to [charmender-spark ChamanderUtils][4]
Note: after resetting the cluster, you won't be able to retrieve the data from previous experiments if that are not saved in an output file.

If you wish, you can further implement the codes you plays with in spark-shell into the analytics piece, please refer to [You can build your own Nessy analytics][5]


[1]: https://influxdb.com/docs/v0.8/introduction/getting_started.html
[2]: http://172.31.1.11:5050/#/
[3]: http://172.31.2.11:31790/#/?host=slave3&hostspec=localhost
[4]: https://github.com/att-innovate/charmander-spark/blob/master/src/main/scala/org/att/charmander/CharmanderUtils.scala
[5]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/analytics/



	