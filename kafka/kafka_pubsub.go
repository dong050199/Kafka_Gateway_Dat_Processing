package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"MQTTClient/config"
	"MQTTClient/model"
)

func Kafka_Pub(Request model.Scada) {
	fmt.Println("save to kafka")
	//jsonString, err := json.Marshal(data)
	//jobString := string(jsonString)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.KAFKA_BOOTSTRAP_SERVERS})
	if err != nil {
		panic(err)
	}
	// Produce messages to topic (asynchronously)
	topic := config.KAFKA_TOPIC_PUB
	var sendkafkadata model.KafkaRwOrder
	sendkafkadata.Method = Request.Command
	sendkafkadata.Id = Request.Order_ID
	data, _ := json.Marshal(sendkafkadata)
	datastring := string(data)
	for _, word := range []string{string(datastring)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
}

func publish(client mqtt.Client, message string, topic string) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
}

func Kafka_Sub(clientMqtt mqtt.Client) {
	fmt.Println("Start receiving from Kafka")
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KAFKA_BOOTSTRAP_SERVERS,
		"group.id":          config.KAFKA_GROUP_ID,
		"auto.offset.reset": config.KAFKA_AUTO_SET_RESET,
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{config.KAFKA_TOPIC_SUB}, nil)

	for {
		msg, err := c.ReadMessage(-1)

		if err == nil {
			fmt.Printf("Received from Kafka %s: %s\n", msg.TopicPartition, string(msg.Value))
			data := string(msg.Value)
			publish(clientMqtt, data, config.MQTT_TOPIC_PUB)

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()
}
