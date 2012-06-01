package mqtt

type MessageType uint8
type ReturnCode uint8
type Header struct{
    messageType MessageType
    dupFlag, retain bool
    qosLevel, length uint8
}
type ConnectFlags struct{
    usernameFlag, passwordFlag, willRetain, willFlag, cleanSession bool
    willQos uint8
}
type Mqtt struct{
    header *Header
    protocolName, topicName, clientIdentifier, willTopic, willMessage, username, password string
    protocolVersion uint8
    connectFlags *ConnectFlags
    keepAliveTimer, messageId uint16
    data []byte
    topics *[]string
    topics_qos *[]uint8
    returnCode ReturnCode
}

const(
    CONNECT = MessageType(iota + 1)
    CONNACK
    PUBLISH
    PUBACK
    PUBREC
    PUBREL
    PUBCOMP
    SUBSCRIBE
    SUBACK
    UNSUBSCRIBE
    UNSUBACK
    PINGREQ
    PINGRESP
    DISCONNECT
)

const(
    ACCEPTED = ReturnCode(iota)
    UNACCEPTABLE_PROTOCOL_VERSION
    IDENTIFIER_REJECTED
    SERVER_UNAVAILABLE
    BAD_USERNAME_OR_PASSWORD
    NOT_AUTHORIZED
)

func getUint8(b []byte, p *int)uint8{
    *p += 1
    return uint8(b[*p-1])
}

func getUint16(b []byte, p *int)uint16{
    *p += 2
    return uint16(b[*p-2] << 8) + uint16(b[*p-1])
}

func getString(b []byte, p *int)string{
    length := int(getUint16(b, p))
    *p += length
    return string(b[*p-length:*p])
}

func getHeader(b []byte, p *int)*Header{
    byte1 := b[*p]
    *p += 1
    header := new(Header)
    header.messageType = MessageType(byte1 & 0xF0 >> 4)
    header.dupFlag = byte1 & 0x08 > 0
    header.qosLevel = uint8(byte1 & 0x06 >> 1)
    header.retain = byte1 & 0x01 > 0
    header.length = getUint8(b, p)
    return header
}

func getConnectFlags(b []byte, p *int)*ConnectFlags{
    bit := b[*p]
    *p += 1
    flags := new(ConnectFlags)
    flags.usernameFlag = bit & 0x80 > 0
    flags.passwordFlag = bit & 0x40 > 0
    flags.willRetain = bit & 0x20 > 0
    flags.willQos = uint8(bit & 0x18 >> 3)
    flags.willFlag = bit & 0x04 > 0
    flags.cleanSession = bit & 0x02 > 0
    return flags
}

func Decode(b []byte)*Mqtt{
    mqtt := new(Mqtt)
    inx := 0
    mqtt.header = getHeader(b, &inx)
    switch mqtt.header.messageType{
        case CONNECT:{
            mqtt.protocolName = getString(b, &inx)
            mqtt.protocolVersion = getUint8(b, &inx)
            mqtt.connectFlags = getConnectFlags(b, &inx)
            mqtt.keepAliveTimer = getUint16(b, &inx)
            mqtt.clientIdentifier = getString(b, &inx)
            if mqtt.connectFlags.willFlag{
                mqtt.willTopic = getString(b, &inx)
                mqtt.willMessage = getString(b, &inx)
            }
            if mqtt.connectFlags.usernameFlag && inx < len(b){
                mqtt.username = getString(b, &inx)
            }
            if mqtt.connectFlags.passwordFlag && inx < len(b){
                mqtt.password = getString(b, &inx)
            }
        }
        case CONNACK:{
            inx += 1
            mqtt.returnCode = ReturnCode(getUint8(b, &inx))
        }
        case PUBLISH:{
            mqtt.topicName = getString(b, &inx)
            if qos := mqtt.header.qosLevel;qos == 1 || qos == 2{
                mqtt.messageId = getUint16(b, &inx)
            }
            mqtt.data = b[inx:len(b)]
            inx = len(b)
        }
        case PUBACK, PUBREC, PUBREL, PUBCOMP, UNSUBACK:{
            mqtt.messageId = getUint16(b, &inx)
        }
        case SUBSCRIBE:{
            if qos := mqtt.header.qosLevel;qos == 1 || qos == 2{
                mqtt.messageId = getUint16(b, &inx)
            }
            topics := make([]string, 0)
            topics_qos := make([]uint8, 0)
            for ; inx < len(b);{
                topics = append(topics, getString(b, &inx))
                topics_qos = append(topics_qos, getUint8(b, &inx))
            }
        }
        case SUBACK:{
            mqtt.messageId = getUint16(b, &inx)
            topics_qos := make([]uint8, 0)
            for ; inx < len(b);{
                topics_qos = append(topics_qos, getUint8(b, &inx))
            }
        }
        case UNSUBSCRIBE:{
            if qos := mqtt.header.qosLevel;qos == 1 || qos == 2{
                mqtt.messageId = getUint16(b, &inx)
            }
            topics := make([]string, 0)
            for ; inx < len(b);{
                topics = append(topics, getString(b, &inx))
            }
        }
        default:
            mqtt = nil
    }
    return mqtt
}
