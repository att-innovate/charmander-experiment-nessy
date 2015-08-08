
// ../spark/bin/spark-submit --class "MaxUsage" --master local[*]  --jars lib/charmander-utils_2.10-1.0.jar target/scala-2.10/max-usage_2.10-1.0.jar

import scala.collection.mutable
import org.apache.spark.{SparkConf, SparkContext}
import org.apache.spark.rdd.RDD

import org.apache.spark.streaming._
import org.json4s.jackson.JsonMethods
import org.att.charmander.CharmanderUtils

//for the delete request
import scalaj.http._


case class MemoryUsage(timestamp: BigDecimal, memory: BigDecimal)

object MaxUsage {

  def main(args: Array[String]) {

    // Create the contexts
    val conf = new SparkConf().setAppName("Charmander-Spark")
    val sc = new SparkContext(conf)
    val ssc = new StreamingContext(sc, Seconds(2))
    //val sqlContext = new org.apache.spark.sql.SQLContext(sc)
    //import sqlContext._

    // Create the queue through which RDDs can be pushed to
    // a QueueInputDStream
    val rddQueue = new mutable.SynchronizedQueue[RDD[List[BigDecimal]]]()

    /*def isDDos(inp: String): Boolean = {
      val allw = inp.split(",").map(wm(_)).sum
      allw > 12
    }*/


    // Create the QueueInputDStream and use it do some processing
    val inputStream = ssc.queueStream(rddQueue)
    inputStream.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println(rdd.first().lift(2))
        //println(Http("http://172.31.1.11:7075/client/task").method("HEAD").asString)
        
        def high_in_bytes(point: List[BigDecimal]): Boolean=
    {
        val in_bytes = BigDecimal(point(2).asInstanceOf[BigInt])
        in_bytes > 20000
    }


        if(rdd.filter(high_in_bytes).count != 0){
          
          println("DDos Attack. Shutting down bonesi-500")

          println(Http("http://172.31.1.11:7075/client/task/bon").method("DELETE").asString)
        }

      }
    })

    ssc.start()


    while (true) {
      //val taskNamesMetered = CharmanderUtils.getMeteredTaskNamesFromRedis()

      //for {taskNameMetered <- taskNamesMetered} {
        rddQueue +=  CharmanderUtils.getRDDForQuery(sc,"charmander-dc","select network_in_bytes from network where interface_name='eth1' AND hostname='slave2' limit 500 ")
      //}

      Thread.sleep(10000)
    }

  }
}