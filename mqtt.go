package mqtt

import ("bytes"
        "errors")

type MessageType uint8
type ReturnCode uint8
type Header struct{
    messageType MessageType
    dupFlag, retain bool
    qosLevel uint8
    length uint32
}
type ConnectFlags struct{
    usernameFlag, passwordFlag, willRetain, willFlag, cleanSession bool
    willQos uint8
}
type Mqtt struct{
    Header *Header
    protocolName, topicName, clientId, willTopic, willMessage, username, password string
    protocolVersion uint8
    ConnectFlags *ConnectFlags
    keepAliveTimer, messageId uint16
    data []byte
    topics []string
    topics_qos []uint8
    ReturnCode ReturnCode
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
    header.length = decodeLength(b, p)
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

func Decode(b []byte)(*Mqtt, error){
    mqtt := new(Mqtt)
    inx := 0
    mqtt.Header = getHeader(b, &inx)
    if mqtt.Header.length != uint32(len(b) - inx){
        return nil, errors.New("Message length is wrong!")
    }
    if msgType := uint8(mqtt.Header.messageType); msgType < 1 || msgType > 14{
        return nil, errors.New("Message Type is invalid!")
    }
    switch mqtt.Header.messageType{
        case CONNECT:{
            mqtt.protocolName = getString(b, &inx)
            mqtt.protocolVersion = getUint8(b, &inx)
            mqtt.ConnectFlags = getConnectFlags(b, &inx)
            mqtt.keepAliveTimer = getUint16(b, &inx)
            mqtt.clientId = getString(b, &inx)
            if mqtt.ConnectFlags.willFlag{
                mqtt.willTopic = getString(b, &inx)
                mqtt.willMessage = getString(b, &inx)
            }
            if mqtt.ConnectFlags.usernameFlag && inx < len(b){
                mqtt.username = getString(b, &inx)
            }
            if mqtt.ConnectFlags.passwordFlag && inx < len(b){
                mqtt.password = getString(b, &inx)
            }
        }
        case CONNACK:{
            inx += 1
            mqtt.ReturnCode = ReturnCode(getUint8(b, &inx))
            if code := uint8(mqtt.ReturnCode);code > 5{
                return nil, errors.New("ReturnCode is invalid!")
            }
        }
        case PUBLISH:{
            mqtt.topicName = getString(b, &inx)
            if qos := mqtt.Header.qosLevel;qos == 1 || qos == 2{
                mqtt.messageId = getUint16(b, &inx)
            }
            mqtt.data = b[inx:len(b)]
            inx = len(b)
        }
        case PUBACK, PUBREC, PUBREL, PUBCOMP, UNSUBACK:{
            mqtt.messageId = getUint16(b, &inx)
        }
        case SUBSCRIBE:{
            if qos := mqtt.Header.qosLevel;qos == 1 || qos == 2{
                mqtt.messageId = getUint16(b, &inx)
            }
            topics := make([]string, 0)
            topics_qos := make([]uint8, 0)
            for ; inx < len(b);{
                topics = append(topics, getString(b, &inx))
                topics_qos = append(topics_qos, getUint8(b, &inx))
            }
            mqtt.topics = topics
            mqtt.topics_qos = topics_qos
        }
        case SUBACK:{
            mqtt.messageId = getUint16(b, &inx)
            topics_qos := make([]uint8, 0)
            for ; inx < len(b);{
                topics_qos = append(topics_qos, getUint8(b, &inx))
            }
            mqtt.topics_qos = topics_qos
        }
        case UNSUBSCRIBE:{
            if qos := mqtt.Header.qosLevel;qos == 1 || qos == 2{
                mqtt.messageId = getUint16(b, &inx)
            }
            topics := make([]string, 0)
            for ; inx < len(b);{
                topics = append(topics, getString(b, &inx))
            }
            mqtt.topics = topics
        }
    }
    return mqtt, nil
}

func setUint8(val uint8, buf *bytes.Buffer){
    buf.WriteByte(byte(val))
}

func setUint16(val uint16, buf *bytes.Buffer){
    buf.WriteByte(byte(val & 0xff00 >> 8))
    buf.WriteByte(byte(val & 0x00ff))
}

func setString(val string, buf *bytes.Buffer){
    length := uint16(len(val))
    setUint16(length, buf)
    buf.WriteString(val)
}

func setHeader(header *Header, buf *bytes.Buffer){
    val := byte(uint8(header.messageType)) << 4
    val |= (boolToByte(header.dupFlag) << 3)
    val |= byte(header.qosLevel) << 1
    val |= boolToByte(header.retain)
    buf.WriteByte(val)
}

func setConnectFlags(flags *ConnectFlags, buf *bytes.Buffer){
    val := boolToByte(flags.usernameFlag) << 7
    val |= boolToByte(flags.passwordFlag) << 6
    val |= boolToByte(flags.willRetain) << 5
    val |= byte(flags.willQos) << 3
    val |= boolToByte(flags.willFlag) << 2
    val |= boolToByte(flags.cleanSession) << 1
    buf.WriteByte(val)
}

func boolToByte(val bool)byte{
    if val{
        return byte(1)
    }
    return byte(0)
}

func Encode(mqtt *Mqtt)([]byte, error){
    err := valid(mqtt)
    if err != nil{
        return nil, err
    }
    var headerbuf, buf bytes.Buffer
    setHeader(mqtt.Header, &headerbuf)
    switch mqtt.Header.messageType{
        case CONNECT:{
            setString(mqtt.protocolName, &buf)
            setUint8(mqtt.protocolVersion, &buf)
            setConnectFlags(mqtt.ConnectFlags, &buf)
            setUint16(mqtt.keepAliveTimer, &buf)
            setString(mqtt.clientId, &buf)
            if mqtt.ConnectFlags.willFlag{
                setString(mqtt.willTopic, &buf)
                setString(mqtt.willMessage, &buf)
            }
            if mqtt.ConnectFlags.usernameFlag && len(mqtt.username) > 0{
                setString(mqtt.username, &buf)
            }
            if mqtt.ConnectFlags.passwordFlag && len(mqtt.password) > 0{
                setString(mqtt.password, &buf)
            }
        }
        case CONNACK:{
            buf.WriteByte(byte(0))
            setUint8(uint8(mqtt.ReturnCode), &buf)
        }
        case PUBLISH:{
            setString(mqtt.topicName, &buf)
            if qos := mqtt.Header.qosLevel;qos == 1 || qos == 2{
                setUint16(mqtt.messageId, &buf)
            }
            buf.Write(mqtt.data)
        }
        case PUBACK, PUBREC, PUBREL, PUBCOMP, UNSUBACK:{
            setUint16(mqtt.messageId, &buf)
        }
        case SUBSCRIBE:{
            if qos := mqtt.Header.qosLevel;qos == 1 || qos == 2{
                setUint16(mqtt.messageId, &buf)
            }
            for i := 0;i < len(mqtt.topics);i += 1{
                setString(mqtt.topics[i], &buf)
                setUint8(mqtt.topics_qos[i], &buf)
            }
        }
        case SUBACK:{
            setUint16(mqtt.messageId, &buf)
            for i := 0;i < len(mqtt.topics_qos);i += 1{
                setUint8(mqtt.topics_qos[i], &buf)
            }
        }
        case UNSUBSCRIBE:{
            if qos := mqtt.Header.qosLevel;qos == 1 || qos == 2{
                setUint16(mqtt.messageId, &buf)
            }
            for i := 0;i < len(mqtt.topics); i += 1{
                setString(mqtt.topics[i], &buf)
            }
        }
    }
    if buf.Len() > 268435455{
        return nil, errors.New("Message is too long!")
    }
    encodeLength(uint32(buf.Len()), &headerbuf)
    headerbuf.Write(buf.Bytes())
    return headerbuf.Bytes(), nil
}

func valid(mqtt *Mqtt)error{
    if msgType := uint8(mqtt.Header.messageType);msgType < 1 || msgType > 14{
        return errors.New("MessageType is invalid!")
    }
    if mqtt.Header.qosLevel > 3 {
        return errors.New("Qos Level is invalid!")
    }
    if mqtt.ConnectFlags != nil && mqtt.ConnectFlags.willQos > 3{
        return errors.New("Will Qos Level is invalid!")
    }
    return nil
}

func decodeLength(b []byte, p *int)uint32{
    m := uint32(1)
    v := uint32(b[*p] & 0x7f)
    *p += 1
    for ; b[*p-1] & 0x80 > 0 ;{
        m *= 128
        v += uint32(b[*p] & 0x7f) * m
        *p += 1
    }
    return v
}

func encodeLength(length uint32, buf *bytes.Buffer){
    if length == 0{
        buf.WriteByte(byte(0))
        return
    }
    var lbuf bytes.Buffer
    for ; length > 0;{
        digit := length % 128
        length = length / 128
        if length > 0{
            digit = digit | 0x80
        }
        lbuf.WriteByte(byte(digit))
    }
    blen := lbuf.Bytes()
    for i := 1;i <= len(blen);i += 1{
        buf.WriteByte(blen[len(blen)-i])
    }
}
