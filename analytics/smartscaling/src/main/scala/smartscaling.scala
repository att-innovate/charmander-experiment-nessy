
import scala.collection.mutable
import org.apache.spark.{SparkConf, SparkContext}
import org.apache.spark.rdd.RDD

import org.apache.spark.streaming._
import org.json4s.jackson.JsonMethods
import org.att.charmander.CharmanderUtils

import scalaj.http._




object SmartScaling {

  val usage = """
    Usage: SmartScaling [--jsonDir dirname] [--low_in num] [--high_in num] [--ratio_inout double]
  """
  val default_low_in = 17000
  val default_high_in = 30000
  val default_ratio_inout = 1.0
  val default_jsonDir = "/files"

  var low_in = default_low_in
  var high_in = default_high_in
  var ratio_inout = default_ratio_inout
  var jsonDir = default_jsonDir

  def configureThresholds(args: Array[String]){
    //val args = args.split(" ")
    def getFlagIndex(flag: String) : Int =  args.indexOf(flag) + 1
    if (args.contains("--jsonDir")) jsonDir = args(getFlagIndex("--jsonDir"))
    if (args.contains("--low_in")) low_in = args(getFlagIndex("--low_in")).toInt
    if (args.contains("--high_in")) high_in = args(getFlagIndex("--high_in")).toInt
    if (args.contains("--ratio_inout")) ratio_inout = args(getFlagIndex("--ratio_inout")).toDouble
    println("Configurations: low_in="+low_in + ", high_in="+high_in+", ratio_inout="+ratio_inout+", jsonDir="+jsonDir)
  }
  
  def isDDOS(rdd: RDD[List[BigDecimal]]): Boolean=
  {
      rdd.filter(w=> w(2).asInstanceOf[BigInt] > high_in).count > 2 && rdd.filter(w=> w(2).asInstanceOf[BigInt] > w(3).asInstanceOf[BigInt]).count > 2   //if in_bytes > out_bytes
  }

  def shouldScaleUp(rdd: RDD[List[BigDecimal]]): Boolean=
  {
      rdd.filter(w=> w(2).asInstanceOf[BigInt] > high_in).count > 2
  }


  def shouldScaleDown(rdd: RDD[List[BigDecimal]]): Boolean=
  {
      rdd.filter(w=> w(2).asInstanceOf[BigInt] < low_in).count >= 2
  }


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
    var ddosRecognized = false 

    inputStream.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println("dns-sl2")
        rdd.foreach(println)

        if(isDDOS(rdd)){
          println("DDos Attack. Shutting down bonesi-500")
          println(Http("http://172.31.1.11:7075/client/task/bon").method("DELETE").asString)
        }

        if(!twoServersUp && shouldScaleUp(rdd)){
          println("Should Scale Up dns-sl3")
          val data= scala.io.Source.fromFile(jsonDir+"/dns-sl3.json").mkString
          println(Http("http://172.31.1.11:7075/client/task").postData(data).header("content-type", "application/json").asString)
          twoServersUp = true
        }

        if(twoServersUp && shouldScaleDown(rdd)){
          println("Should Scale Down and shut down dns-sl3")
          println(Http("http://172.31.1.11:7075/client/task/dns-sl3").method("DELETE").asString)
          twoServersUp = false
        }

      }
    })

    ssc.start()

    while (true) {
      val (col, rdd) =  CharmanderUtils.getRDDAndColumnsForQuery(sc,CharmanderUtils.VECTOR_DB,"select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave2' limit 5")
      rddQueue += rdd

      Thread.sleep(5000)
    }

   
  }

}