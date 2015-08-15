name := "smartscaling"

version := "1.0"

scalaVersion := "2.10.4"

libraryDependencies += "org.apache.spark" %% "spark-core" % "1.2.0" % "provided"

libraryDependencies += "org.apache.spark" %% "spark-streaming" % "1.2.0"  % "provided"

libraryDependencies +=  "org.scalaj" %% "scalaj-http" % "1.1.5"