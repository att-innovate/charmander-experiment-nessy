
import scala.collection.mutable
import org.apache.spark.{SparkConf, SparkContext}
import org.apache.spark.rdd.RDD

import org.apache.spark.streaming._
import org.json4s.jackson.JsonMethods
import org.att.charmander.CharmanderUtils

import scalaj.http._


/* SmartScaling
 * -------------------
 * SmartScaling is the intelligence piece of Project Nessy. It streams rdd data from the DNS 
 * servers and runs analysis on the rdd data to identify whether the system needs to scale up,
 * scale down, or react to the DDoS attack.
 *
 * To edit the json path for the container runnign SmartScaling, you can specify --jsonDir
 *
 * To change the auto-scaling policy and DDoS detection, you can specify these following thresholds:
 *  --high_in: if we observe network_in_byte to be larger than high_in, we call the schdeule to spin up
 *            another server
 *  --low_in:  if we observe network_in_byte to be less than low_in, and if there are 2 DNS servers running,
 *            we call the scheduler to shut down the second DNS server
 *  --ratio_inout: if we observe network_in_byte/net_out_byte to be larger than ratio_inout, and if
 *            the network_in_byte larger than high_in, we regard it as DDoS attack and kill it
 *  --baseline: if network_in_byte is lower than baseline, we don't do anything
 *
 */

object SmartScaling {
  //Constants and Initial values
  val usage = """
    Usage: SmartScaling [--jsonDir dirname] [--low_in num] [--high_in num] [--baseline num] [--ratio_inout double]
  """
  val shedulerAPI_addr = "http://172.31.1.11:7075/client/task"
  val shedulerAPI_bonesi_addr = shedulerAPI_addr + "/bonesi"
  val shedulerAPI_dns3_addr = shedulerAPI_addr + "/dns-sl3"
  val num_tolerate = 2

  val default_low_in = 17000
  val default_high_in = 30000
  val default_ratio_inout = 0.6
  val default_jsonDir = "/files"
  val default_baseline = 10000

  var low_in = default_low_in
  var high_in = default_high_in
  var ratio_inout = default_ratio_inout
  var jsonDir = default_jsonDir
  var baseline = default_baseline

  //Configure the user inputs 
  def configureThresholds(args: Array[String]){
    def getFlagIndex(flag: String) : Int =  args.indexOf(flag) + 1
    if (args.contains("--jsonDir")) jsonDir = args(getFlagIndex("--jsonDir"))
    if (args.contains("--low_in")) low_in = args(getFlagIndex("--low_in")).toInt
    if (args.contains("--high_in")) high_in = args(getFlagIndex("--high_in")).toInt
    if (args.contains("--baseline")) high_in = args(getFlagIndex("--baseline")).toInt
    if (args.contains("--ratio_inout")) ratio_inout = args(getFlagIndex("--ratio_inout")).toDouble
    println("Configurations: low_in="+low_in + ", high_in="+high_in + ", baseline="+baseline+", ratio_inout="+ratio_inout+", jsonDir="+jsonDir)
  }
  

  //Algorithms for Auto-scaling and DDoS detection
  def network_in (value: List[BigDecimal]) : BigInt = value(2).asInstanceOf[BigInt]
  def network_out (value: List[BigDecimal]) : BigInt = value(3).asInstanceOf[BigInt]

  def isDDOS(rdd: RDD[List[BigDecimal]]): Boolean=
  {
      shouldScaleUp(rdd) && (rdd.filter(value => (network_in(value).toDouble / network_out(value).toDouble) > ratio_inout).count > num_tolerate)
  }

  def shouldScaleUp(rdd: RDD[List[BigDecimal]]): Boolean=
  {
      rdd.filter(value => network_in(value) > high_in).count > num_tolerate
  }

  def shouldScaleDown(rdd: RDD[List[BigDecimal]]): Boolean=
  {
      rdd.filter(value => network_in(value) < low_in && network_in(value) > baseline).count > num_tolerate
  }


  //Commands sent to schedulers
  def killDDoS () { 
    println(Http(shedulerAPI_bonesi_addr).method("DELETE").asString)
  }
  
  def startNewDNS () {
    val data= scala.io.Source.fromFile(jsonDir+"/dns-sl3.json").mkString
    println(Http(shedulerAPI_addr).postData(data).header("content-type", "application/json").asString)           
  }

  def shutDownDNS () {
    println(Http(shedulerAPI_dns3_addr).method("DELETE").asString)
  }

  // Main workflow: streaming data 
  def main(args: Array[String]) {
     //Parse arguments, set as default values if user does not define it
    if (args.length == 0) println(usage)
    else configureThresholds(args)
    
    //Configure Spark Contents
    val conf = new SparkConf().setAppName("Charmander-Nessy")
    val sc = new SparkContext(conf)
    val ssc = new StreamingContext(sc, Seconds(5))

    // Create the queue through which RDDs can be pushed to
    // a QueueInputDStream
    val rddQueue = new mutable.SynchronizedQueue[RDD[List[BigDecimal]]]()
    // Create the QueueInputDStream and use it do some processing
    val inputStream = ssc.queueStream(rddQueue)
    var twoServersUp = false 

    // Constantly analyzing streamed data, recognize situations and send relative commands to scheduler
    inputStream.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println("dns-sl2")
        rdd.foreach(println)

        if(isDDOS(rdd)){
          println("DDos Attack. Shutting down bonesi-500")
          killDDoS()
        }

        if(!twoServersUp && shouldScaleUp(rdd)){
          println("Should Scale Up dns-sl3")
          startNewDNS()
          twoServersUp = true
        }

        if(twoServersUp && shouldScaleDown(rdd)){
          println("Should Scale Down and shut down dns-sl3")
          shutDownDNS()
          twoServersUp = false
        }

      }
    })

    ssc.start()

    // Collect data from InfluxDB every 5 seconds and add them into  data stream
    while (true) {
      val (col, rdd) =  CharmanderUtils.getRDDAndColumnsForQuery(sc,CharmanderUtils.VECTOR_DB,"select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave2' limit 5")
      rddQueue += rdd
      Thread.sleep(5000)
    } 
  }
}