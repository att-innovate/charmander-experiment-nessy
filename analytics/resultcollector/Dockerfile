FROM phusion/baseimage:0.9.17

ENV HOME /root
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -qq
RUN apt-get -y install git wget

# Install Java
ENV JAVA_HOME /usr/jdk1.8.0_25

RUN cd /tmp \
    && wget \
    --no-check-certificate \
    --no-cookies \
    --header "Cookie: oraclelicense=accept-securebackup-cookie" \
    --progress=bar:force \
    download.oracle.com/otn-pub/java/jdk/8u25-b17/jdk-8u25-linux-x64.tar.gz \
    && tar -xvf jdk-8u25-linux-x64.tar.gz -C /usr/ \
    && ln -s $JAVA_HOME /usr/java \
    && rm -rf $JAVA_HOME/src.zip $JAVA_HOME/javafx-src.zip $JAVA_HOME/man

ENV PATH $PATH:$JAVA_HOME/bin


# Install sbt
RUN cd /tmp \
    && wget \
    --no-check-certificate \
    --no-cookies \
    --progress=bar:force \
    https://dl.bintray.com/sbt/native-packages/sbt/0.13.7/sbt-0.13.7.tgz \
    && tar -xvf sbt-0.13.7.tgz -C /usr/ \
    && /usr/sbt/bin/sbt -batch

ENV PATH $PATH:/usr/sbt/bin


# Install resultcollector
ENV RESULT_COLLECTOR_DIR /resultcollector

RUN mkdir -p $RESULT_COLLECTOR_DIR/src
RUN mkdir -p $RESULT_COLLECTOR_DIR/project

ADD src/ $RESULT_COLLECTOR_DIR/src

COPY build.sbt $RESULT_COLLECTOR_DIR/
COPY project/assembly.sbt $RESULT_COLLECTOR_DIR/project/

RUN cd $RESULT_COLLECTOR_DIR && sbt assembly


# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Execute ResultCollector
ENTRYPOINT ["java", "-jar", "/resultcollector/target/scala-2.10/resultcollector-assembly-1.0.jar"]



