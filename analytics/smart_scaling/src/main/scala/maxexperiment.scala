
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

    val high_ave_in = 28000
    val high_ave_out = 63000

    var twoServersUp = false

    // Create the queue through which RDDs can be pushed to
    // a QueueInputDStream
    val rddQueue_1 = new mutable.SynchronizedQueue[RDD[List[BigDecimal]]]()
    val rddQueue_2 = new mutable.SynchronizedQueue[RDD[List[BigDecimal]]]()

    // Create the QueueInputDStream and use it do some processing
    val inputStream_1 = ssc.queueStream(rddQueue_1)
    inputStream_1.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println("dns-sl2:" + rdd.first())

        def isDDOS(rdd: RDD[List[BigDecimal]]): Boolean=
    {
        rdd.filter(w=> w(2).asInstanceOf[BigInt] > high_ave_in).count > 30 && rdd.filter(w=> w(2).asInstanceOf[BigInt] > w(3).asInstanceOf[BigInt]).count > 20   //if in_bytes > out_bytes
    }
 
        def shouldScaleUp(rdd: RDD[List[BigDecimal]]): Boolean=
    {
        rdd.filter(w=> w(2).asInstanceOf[BigInt] > high_ave_in).count > 40 && rdd.filter(w=>w(3).asInstanceOf[BigInt] > high_ave_out).count > 40
    }

        def shouldScaleDown(rdd: RDD[List[BigDecimal]]): Boolean=
    {
        rdd.filter(w=> w(2).asInstanceOf[BigInt] < high_ave_in).count > 60 && rdd.filter(w=>w(3).asInstanceOf[BigInt] < high_ave_out).count > 60
    }

        if(isDDOS(rdd)){  //Naive DDos detection: in_bytes>out_bytes for more than 5 cases
          println("DDos Attack. Shutting down bonesi-500")
          println(Http("http://172.31.1.11:7075/client/task/bon").method("DELETE").asString)
        }


         if(!twoServersUp && shouldScaleUp(rdd)){  //Naive DDos detection: in_bytes>out_bytes for more than 5 cases
          println("Should Scale Up dns-sl3")
          val data= scala.io.Source.fromFile("/Users/Karen/workArea/charmander/experiments/nessy/services/dns/dns-sl3.json").mkString 
          println(Http("http://172.31.1.11:7075/client/task").postData(data).header("content-type", "application/json").asString)
          twoServersUp = true
        }

        if(twoServersUp && shouldScaleDown(rdd)){  //Naive DDos detection: in_bytes>out_bytes for more than 5 cases
          println("Should Scale Down and shut down dns-sl3")
          println(Http("http://172.31.1.11:7075/client/task/dns-sl3").method("DELETE").asString)
          twoServersUp = false
        }

      }
    })

    
    val inputStream_2 = ssc.queueStream(rddQueue_2)
    inputStream_2.foreachRDD(rdd => {
      if (rdd.count != 0) {
        println("dns-sl3:" + rdd.first())        
      }
    })


    ssc.start()


    while (true) {
      val (col, rdd) =  CharmanderUtils.getRDDAndColumnsForQuery(sc,"charmander-dc","select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave2' limit 100 ")
      //val (index_time,index_inbytes,index_outbytes) = (col.indexOf("time"), col.indexOf("network_in_bytes"), col.indexOf("network_out_bytes"))
      //network_rddQueue += rdd.map(rdd => List(BigDecimal(rdd(index_time).asInstanceOf[BigInt]),BigDecimal(rdd(index_inbytes).asInstanceOf[BigInt]),BigDecimal(rdd(index_outbytes).asInstanceOf[BigInt])))
      rddQueue_1 += rdd

      val (col_2, rdd_2) =  CharmanderUtils.getRDDAndColumnsForQuery(sc,"charmander-dc","select network_in_bytes,network_out_bytes from network where interface_name='eth1' AND hostname='slave3' limit 100 ")
      rddQueue_2 += rdd_2
      /*val (col_2,rdd_2) = CharmanderUtils.getRDDAndColumnsForQuery(sc,"charmander-dc","select * from stats where hostname='slave2' limit 100 ")
      val (index_time_2,index_cpu_user,index_cpu_sys,index_mem) = (col_2.indexOf("time"), col_2.indexOf("cpu_usage_user"), col_2.indexOf("cpu_usage_system"), col_2.indexOf("memory_usage"))
      stats_rddQueue += rdd_2.map(rdd_2 => List(BigDecimal(rdd_2(index_time_2).asInstanceOf[BigInt]),BigDecimal(rdd_2(index_mem).asInstanceOf[BigInt]),BigDecimal(rdd_2(index_cpu_user).asInstanceOf[BigInt]),BigDecimal(rdd_2(index_cpu_sys).asInstanceOf[BigInt])))
      */

      Thread.sleep(10000)


    }

  }
}