Run a Nessy Script
------------------


These scripts will automatically set up everything and run some normalloads and DDos attacks to the DNS server. It does not include smartscaling.

The script will output a logfile.csv. It includes the number of users and whether there is a DDoS attacks(bonesi), and the timestamps showing the starting time of a simulation that causes the change.it will appear in your local charmadner folder. They can be a reference if you wish to run analysis based on these information.

In your Charmander folder:

    ./experiments/nessy/bin/script_labels_30min

You will need to wait for some time to ensure all containers are up and running, and the data_collector is working.

Then, open the [Mesos](http://172.31.1.11:5050/#/) to check tasks that are actively running. It should look like this:

![image](https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/MesosExp.png?raw=true)


If you wish to see the data generated on run-time, open the [Vector](http://http://172.31.2.11:31790/#/?host=slave2&hostspec=localhost) to observe the metrics collected from DNS server on slave2.


When the script complets, you can collect all the data it generated from InfluxDB into a csv file by querying from Spark. (InfluxDB itself does not provide csv conversion)
In your SPARK_FOLDER, run:

	