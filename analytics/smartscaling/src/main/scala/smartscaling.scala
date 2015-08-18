
import scala.collection.mutable
import org.apache.spark.{SparkConf, SparkContext}
import org.apache.spark.rdd.RDD

import org.apache.spark.streaming._
import org.json4s.jackson.JsonMethods
import org.att.charmander.CharmanderUtils

import scalaj.http._




object SmartScaling {

  def main(args: Array[String]) {

    // set the jsonDir, the directory that contains the json files
    // for starting containers
    var jsonDir = "/files"

    if (args.length == 2 && args(0) == "--jsonDir") {
        jsonDir = args(1)
    }

    println("jsonDir: " + jsonDir)

    // Create the contexts
    val conf = new SparkConf().setAppName("Charmander-Nessy")
    val sc = new SparkContext(conf)
    val ssc = new StreamingContext(sc, Seconds(5))

    val high_ave_in = 30000
    val high_ave_out = 60000

    val low_ave_in = 15000

    var twoServersUp = false
    var ddosRecognized = false

    // Create the queue through which RDDs can be pushed to
    // a QueueInputDStream
    val rddQueue = new mutable.SynchronizedQueue[RDD[List[BigDecimal]]]()

    // Create the QueueInputDStream and use it do some processing
    val inputStream_1 = ssc.queueStream(rddQueue)
    inputStream_1.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println("dns-sl2")
        rdd.foreach(println)

        def isDDOS(rdd: RDD[List[BigDecimal]]): Boolean=
        {
            rdd.filter(w=> w(2).asInstanceOf[BigInt] > high_ave_in).count > 2 && rdd.filter(w=> w(2).asInstanceOf[BigInt] > w(3).asInstanceOf[BigInt]).count > 2   //if in_bytes > out_bytes
        }
 
        def shouldScaleUp(rdd: RDD[List[BigDecimal]]): Boolean=
        {
            rdd.filter(w=> w(2).asInstanceOf[BigInt] > high_ave_in).count > 2 && rdd.filter(w=>w(3).asInstanceOf[BigInt] > high_ave_out).count > 2
        }

        def shouldScaleDown(rdd: RDD[List[BigDecimal]]): Boolean=
        {
            rdd.filter(w=> w(2).asInstanceOf[BigInt] < low_ave_in).count >= 2
        }

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
      val (col, rdd) =  CharmanderUtils.getRDDAndColumnsForQuery(sc,CharmanderUtils.VECTOR_DB,"select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave2' limit 5 ")
      rddQueue += rdd

      Thread.sleep(10000)
    }

  }
}