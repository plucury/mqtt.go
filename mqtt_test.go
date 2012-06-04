package mqtt

import (
    "testing"
    "fmt"
    )

var bitCnt = uint32(0)

func Test(t *testing.T){
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

func initTest()*Mqtt{
    mqtt := new(Mqtt)
    mqtt.header = new(Header)
    mqtt.header.messageType = MessageType(1)
    mqtt.protocolName = "MQIsdp"
    mqtt.protocolVersion = uint8(3)
    mqtt.connectFlags = new(ConnectFlags)
    mqtt.connectFlags.usernameFlag = true
    mqtt.connectFlags.passwordFlag = true
    mqtt.connectFlags.willRetain = false
    mqtt.connectFlags.willQos = uint8(1)
    mqtt.connectFlags.willFlag = true
    mqtt.connectFlags.cleanSession = true
    mqtt.keepAliveTimer = uint16(10)
    mqtt.clientId = "xixihaha"
    mqtt.willTopic = "topic"
    mqtt.willMessage = "message"
    mqtt.username = "name"
    mqtt.password = "pwd"
    return mqtt
}

func printByte(b byte){
    bitCnt += 1
    out := make([]uint8, 8)
    val := uint8(b)
    for i := 1; val > 0;i += 1{
        foo := val % 2
        val = (val - foo) / 2
        out[8 - i] = foo
    }
    fmt.Println(bitCnt, out)
}

func printBytes(b []byte){
    for i := 0;i < len(b);i += 1{
        printByte(b[i])
    }
}

func printMqtt(mqtt *Mqtt){
    fmt.Printf("MQTT = %+v\n", *mqtt)
    fmt.Printf("Header = %+v\n", *mqtt.header)
    if mqtt.connectFlags != nil{
        fmt.Printf("ConnectFlags = %+v\n", *mqtt.connectFlags)
    }
}
