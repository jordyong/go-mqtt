package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func ConnectMQTT(clientName, brokerURL string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID(clientName)
	opts.SetDefaultPublishHandler(f)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Failed to init MQTT: %s\n", token.Error())
		return nil
	}
	PublishMQTT(client, "test/topic", "test")
	return client
}

func SubscribeMQTT(c mqtt.Client, topic string) error {
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func PublishMQTT(c mqtt.Client, topic, msg string) error {
	token := c.Publish(topic, 0, false, msg)
	token.Wait()
	return nil
}
