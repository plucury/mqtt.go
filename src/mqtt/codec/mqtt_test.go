package codex

import (
	"fmt"
	"testing"
)

var bitCnt = uint32(0)

func Test(t *testing.T) {
	mqtt := initTest()
	fmt.Println("------ Origin MQTT Object")
	printMqtt(mqtt)
	fmt.Println("------ Encode To Binary")
	bits, _ := Encode(mqtt)
	printBytes(bits)
	fmt.Println("------ Decode To Object")
	newMqtt, _ := Decode(bits)
	printMqtt(newMqtt)
}

func initTest() *Mqtt {
	mqtt := new(Mqtt)
	mqtt.Header = new(Header)
	mqtt.Header.MessageType = MessageType(1)
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
	return mqtt
}

func printByte(b byte) {
	bitCnt += 1
	out := make([]uint8, 8)
	val := uint8(b)
	for i := 1; val > 0; i += 1 {
		foo := val % 2
		val = (val - foo) / 2
		out[8-i] = foo
	}
	fmt.Println(bitCnt, out)
}

func printBytes(b []byte) {
	for i := 0; i < len(b); i += 1 {
		printByte(b[i])
	}
}

func printMqtt(mqtt *Mqtt) {
	fmt.Printf("MQTT = %+v\n", *mqtt)
	fmt.Printf("Header = %+v\n", *mqtt.Header)
	if mqtt.ConnectFlags != nil {
		fmt.Printf("ConnectFlags = %+v\n", *mqtt.ConnectFlags)
	}
}
