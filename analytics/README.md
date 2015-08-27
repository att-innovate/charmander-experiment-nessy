Nessy Intelligence: Analytics
-----------------------------

Analytics is the intelligence piece of Project Nessy. Smartscaling, the template of analytics we provide, streams rdd data from the DNS 
servers and runs analysis on the rdd data for auto-scaling and DDoS detection. Here we have provided a simple implementation of auto-scaling and DDoS detection. It can be easily plugined with your own implementation.

To develop your own auto-scaling policy or DDoS detection, follow the following instructions:


1. Copy the whole smartscaling folder into your spark folder

2. Installed the sbt assembly (Please refer to https://github.com/sbt/sbt-assembly)

3. To start a interactive shell to play around with the scalaj and CharmanderUtils:

		../bin/spark-shell --jars $(echo ./lib/*.jar | tr ' ' ',')

   Inside the spark-shell:
   		>>> import scalaj.http._
   		>>> import org.att.charmander.CharmanderUtils
   		>>> ......(your own code)

4. Then you might wish to make changes in the smartscaling script. 
		
		cd ./src/main/scala

	Open smartscaling.scala in your favorite editor. 
	You can change the query command to get other metrics from InfluxDB. 
	You can change the sleeptime to set the interval for data streaming.
	You can replace the isDDoS, shouldScaleUp, shouldScaleDown fuctions with your own.

5. To compile your own smartscaling.scala:

	sbt assembly

6.	To run the program, you should set the -jsonDir to be your own directory that stores the json file.
	Optionally, you can also set up the threshold for low_in (lowest network_in_byte), high_in, ratio_inout(ratio of network_in_byte over network_outbyte) and baseline (below which there won't be any scaling up/down on the dns server):

	./bin/spark-submit --class "SmartScaling" --master local[*]  $YOUR_SPARK_FOLDER/smartscaling/target/scala-2.10/smartscaling-assembly-1.0.jar --jsonDir $YOUR_SPARK_FOLDER/files/  [optional: --high_in --low_in --baseline]
