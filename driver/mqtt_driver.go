package driver

import (
	"MQTTClient/config"
	"MQTTClient/kafka"
	"MQTTClient/model"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	//run controller here
	var datareceive model.Scada
	err := json.Unmarshal(msg.Payload(), &datareceive)
	if err != nil {
		//panic(err)
		fmt.Println(err)
		return
	}
	kafka.Kafka_Pub(datareceive)

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
var ClientMQ mqtt.Client

func Connect() {
	var broker = config.MQTT_BROKER_HOST
	var port = config.MQTT_PORT
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(config.MQTT_CLIENT_ID)
	opts.SetUsername(config.MQTT_USERNAME)
	opts.SetPassword(config.MQTT_PASSWORD)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	sub(client)
	ClientMQ = client
}

func sub(client mqtt.Client) {
	topic := config.MQTT_TOPIC_SUB
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
