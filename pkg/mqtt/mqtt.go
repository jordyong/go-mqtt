package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttService struct {
	client mqtt.Client
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func NewMqttService(clientName, brokerURL string) (*MqttService, error) {
	fmt.Printf("Connecting to %s at %s\n", clientName, brokerURL)
	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID(clientName)
	opts.SetDefaultPublishHandler(f)
	c := mqtt.NewClient(opts)

	return &MqttService{client: c}, nil
}

func (ms *MqttService) IsConnected() bool {
	return ms.client.IsConnected()
}

func (ms *MqttService) Connect() error {
	if token := ms.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (ms *MqttService) Disconnect() error {
	ms.client.Disconnect(250)
	return nil
}

func (ms *MqttService) Subscribe(topic string, handler mqtt.MessageHandler) error {
	if token := ms.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (ms *MqttService) Publish(topic string, msg any) error {
	token := ms.client.Publish(topic, 0, false, msg)
	token.Wait()
	return nil
}
