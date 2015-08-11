
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

    // Create the queue through which RDDs can be pushed to
    // a QueueInputDStream
    val network_rddQueue = new mutable.SynchronizedQueue[RDD[List[BigDecimal]]]()
    val stats_rddQueue = new mutable.SynchronizedQueue[RDD[List[BigDecimal]]]()

    // Create the QueueInputDStream and use it do some processing
    val network_inputStream = ssc.queueStream(network_rddQueue)
    network_inputStream.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println("network:" + rdd.first())

        def in_higher_than_out(point: List[BigDecimal]): Boolean=
    {
        val in_bytes = point(1)
        val out_bytes = point(2)
        in_bytes > out_bytes
    }
 
        if(rdd.filter(in_higher_than_out).count > 5){  //Naive DDos detection: in_bytes>out_bytes for more than 5 cases
          println("DDos Attack. Shutting down bonesi-500")
          println(Http("http://172.31.1.11:7075/client/task/bon").method("DELETE").asString)
        }

      }
    })

    
    val stats_inputStream = ssc.queueStream(stats_rddQueue)
    stats_inputStream.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println("stats:" + rdd.first())        
      }
    })


    ssc.start()


    while (true) {
      val (col, rdd) =  CharmanderUtils.getRDDAndColumnsForQuery(sc,"charmander-dc","select * from network where interface_name='eth1' AND hostname='slave2' limit 100 ")
      val (index_time,index_inbytes,index_outbytes) = (col.indexOf("time"), col.indexOf("network_in_bytes"), col.indexOf("network_out_bytes"))
      network_rddQueue += rdd.map(rdd => List(BigDecimal(rdd(index_time).asInstanceOf[BigInt]),BigDecimal(rdd(index_inbytes).asInstanceOf[BigInt]),BigDecimal(rdd(index_outbytes).asInstanceOf[BigInt])))
      
      val (col_2,rdd_2) = CharmanderUtils.getRDDAndColumnsForQuery(sc,"charmander-dc","select * from stats where hostname='slave2' limit 100 ")
      val (index_time_2,index_cpu_user,index_cpu_sys,index_mem) = (col_2.indexOf("time"), col_2.indexOf("cpu_usage_user"), col_2.indexOf("cpu_usage_system"), col_2.indexOf("memory_usage"))
      stats_rddQueue += rdd_2.map(rdd_2 => List(BigDecimal(rdd_2(index_time_2).asInstanceOf[BigInt]),BigDecimal(rdd_2(index_mem).asInstanceOf[BigInt])))


      Thread.sleep(10000)


    }

  }
}