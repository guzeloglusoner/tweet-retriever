package org.apache.spark.examples.streaming.twitter

import java.util.Properties
import org.apache.spark.SparkConf
import org.apache.spark.streaming.twitter.TwitterUtils
import org.apache.spark.streaming.{Seconds, StreamingContext}
import twitter4j.{FilterQuery, Status}
import scala.io.Source
import net.liftweb.json._
import org.apache.kafka.clients.producer.{KafkaProducer, ProducerRecord}
import org.apache.spark.streaming.dstream.{DStream, ReceiverInputDStream}

object TwitterLocations {
  implicit val formats: DefaultFormats.type = net.liftweb.json.DefaultFormats

  case class TwitterCredentials(consumerKey: String, consumerSecret: String, accessToken: String, accessTokenSecret: String)

  def main(args: Array[String]): Unit = {

    /* Read the file */
    var str = ""
    for (line <- Source.fromFile("C:\\DEV\\tweet-retriever\\server\\src\\resources\\env.json").getLines) str += line

    /* Map the JSON file into a case class */
    val twitterCredentials = parse(str).json.extract[TwitterCredentials]

    /* Set the system properties so that Twitter4j library used by twitter stream */
    setTwitterCredentials(twitterCredentials)

    /* Get bounding boxes of locations for which to retrieve Tweets from command line */
    val boundingBoxes = {
      val southWest = Array(25.62128903, 35.80768033)
      val northEast = Array(44.81766374, 42.29699998)

      Array(southWest, northEast)
    }
    val sparkConf = new SparkConf().setAppName("TwitterLocations")
    /* Sets spark master to local */
    if (!sparkConf.contains("spark.master")) {
      sparkConf.setMaster("local[*]")
    }
    val ssc: StreamingContext = new StreamingContext(sparkConf, Seconds(1))
    ssc.sparkContext.setLogLevel("INFO")
    val locationsQuery: FilterQuery = new FilterQuery().locations(boundingBoxes: _*)

    /* Print the tweets with given coordinates*/
    val tweets = TwitterUtils.createFilteredStream(ssc, None, Some(locationsQuery))
    /*.map(tweet => {
      tweet.getCreatedAt
      val latitude = Option(tweet.getGeoLocation).map(l => s"${l.getLatitude},${l.getLongitude}")
      val place = Option(tweet.getPlace).map(_.getName)
      val location = latitude.getOrElse(place.getOrElse("(no location)"))
      val text = tweet.getText.replace('\n', ' ').replace('\r', ' ').replace('\t', ' ')
      val user: Option[String] = Option(tweet.getUser).map(u => s"${u.getName}")
      s"Username: ${user.getOrElse("testUser")}\n Location: $location\n Tweet:$text\n"
    })*/

    val streamedTweets: DStream[(String, String, String, String)] = tweets.map(t => (t.getText, t.getUser.getName, t.getUser.getScreenName, t.getCreatedAt.toString))

    streamedTweets.foreachRDD { (rdd, time) =>

      rdd.foreachPartition {
        partitionIter =>
          val props = new Properties()
          val bootstrap = "10.65.135.159:29092" //-- your external ip
        val zooKeeper = "10.65.135.159:2181"
          props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer")
          props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer")
          props.put("bootstrap.servers", bootstrap)
          props.put("zookeeper.connect", zooKeeper)
          val producer = new KafkaProducer[String, String](props)
          partitionIter.foreach { elem =>
            val dat: String = elem._1
            val data = new ProducerRecord[String, String]("tweet", null, dat) // name of Kafka topic
            producer.send(data)
          }
          producer.flush()
          producer.close()
      }
    }
    ssc.start()
    ssc.awaitTermination()
  }

  private def setTwitterCredentials(twitterCredentials: TwitterCredentials) = {
    System.setProperty("twitter4j.oauth.consumerKey", twitterCredentials.consumerKey)
    System.setProperty("twitter4j.oauth.consumerSecret", twitterCredentials.consumerSecret)
    System.setProperty("twitter4j.oauth.accessToken", twitterCredentials.accessToken)
    System.setProperty("twitter4j.oauth.accessTokenSecret", twitterCredentials.accessTokenSecret)
  }

}