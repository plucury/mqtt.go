package client

import (
	"errors"
	. "mqtt/codec"
)

type MqttClient struct {
	ClientId, Server string
}

func (mc MqttClient) connect() {
	mqtt := new(Mqtt)
	mqtt.Header = new(Header)
	mqtt.Header.MessageType = CONNECT
	mqtt.ProtocolName = "MQIsdp"
	mqtt.ProtocolVersion = uint8(3)
	mqtt.ConnectFlags = new(ConnectFlags)
	mqtt.ConnectFlags.UsernameFlag = true
	mqtt.ConnectFlags.PasswordFlag = true
	mqtt.ConnectFlags.WillRetain = false
	mqtt.ConnectFlags.WillQos = uint8(1)
	mqtt.ConnectFlags.WillFlag = true
	mqtt.ConnectFlags.CleanSession = true
	mqtt.KeepAliveTimer = uint16(10)
	mqtt.ClientId = "xixihaha"
	mqtt.WillTopic = "topic"
	mqtt.WillMessage = "message"
	mqtt.Username = "name"
	mqtt.Password = "pwd"
}
