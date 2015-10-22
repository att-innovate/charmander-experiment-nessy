// The MIT License (MIT)
//
// Copyright (c) 2015 AT&T
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import java.io.{File, PrintWriter, IOException}
import java.util.Calendar
import java.text.SimpleDateFormat
import org.json4s.jackson.JsonMethods


object ResultCollector {

  val usage =
    """
    Usage: java -jar resultcollector-assembly-1.0.jar -db "database" -query "InfluxDBQuery" -name "NameForTheSavedResult" [-dir "OutputDirectory"]
    """

  val INFLUXDB_HOST = "172.31.2.11"
  val INFLUXDB_PORT = 31410


  var database = "foo"
  var query = "foo"
  var resultname = "foo"
  var dir = "."


  def handleCmdlineArguments(args: Array[String]) = {
    def getFlagIndex(flag: String) : Int =  args.indexOf(flag) + 1

    if (args.contains("-db")) database = args(getFlagIndex("-db"))
    if (args.contains("-query")) query = args(getFlagIndex("-query"))
    if (args.contains("-name")) resultname = args(getFlagIndex("-name"))
    if (args.contains("-dir")) dir = args(getFlagIndex("-dir"))

    println("Cmdline Arguments: db: " + database + ", query: " + query + ", resultname: " + resultname)
  }

  def sendQueryToInfluxDB(databaseName: String, query: String): String = try {
    val in = scala.io.Source.fromURL("http://"
      + INFLUXDB_HOST
      + ":"
      + INFLUXDB_PORT
      + "/db/"+databaseName+"/series?u=root&p=root&q="
      + java.net.URLEncoder.encode(query),
      "utf-8")
    var data = ""
    for (line <- in.getLines)
      data = line
    data

  } catch {
    case e: IOException => ""
    case e: java.net.ConnectException => ""
  }

  def exportDataToCSV(databaseName: String, sqlQuery: String, fileName: String) = {
    val writer = new PrintWriter(new File(fileName))

    val rawData = sendQueryToInfluxDB(databaseName,sqlQuery)
    if (rawData.length != 0) {
      val json = JsonMethods.parse(rawData)
      val columns = (json \\ "columns").values
      val points  = (json \\ "points").values

      // handle header
      var header = ""
      for (columnName <- columns.asInstanceOf[List[String]]) {
        if (header != "") header = header + ","

        header = header + columnName
      }

      writer.println(header)

      for (point <- points.asInstanceOf[List[Any]]) {
        var line = ""
        for (elmt <- point.asInstanceOf[List[Any]]) {
          if (line != "") line = line + ","

          line = line + elmt.toString
        }
        writer.println(line)
      }

    }

    writer.close()
  }

  def getTimestamp(): String = {
    val now = Calendar.getInstance().getTime()
    val dataFormat = new SimpleDateFormat("yyyyMMddkkmm")
    dataFormat.format(now)
  }

  def main(args: Array[String]) {
    if (args.length < 6) println(usage)
    else handleCmdlineArguments(args)

    val outputFilePath = dir + "/" + getTimestamp() + "-"+resultname + ".csv"
    exportDataToCSV(database, query, outputFilePath)
  }

}