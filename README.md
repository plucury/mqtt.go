#mqtt.go

an MQTT encoder & decoder,written in GO lang (*version 1.0.1*)

##Functions
* `func Encode(mqtt *Mqtt) (byte[], error)`

	Convert Mqtt struct to bit stream.

* `func Decode(bitstream byte[]) (*Mqtt, error)`

	Convert bit stream to Mqtt struct.
	


