package main

import "fmt"
import "net"
import "encoding/binary"

type DNSHeader struct {
    ID uint16
    Packed uint16
    QDcount uint16
    ANcount uint16
    NScount uint16
    ARcount uint16
}

func NewDNSHeaderFromBytes(buf []byte) *DNSHeader {
    id := binary.BigEndian.Uint16(buf[0:])
    pack := binary.BigEndian.Uint16(buf[2:])
    qd := binary.BigEndian.Uint16(buf[4:])
    an := binary.BigEndian.Uint16(buf[6:])
    ns := binary.BigEndian.Uint16(buf[8:])
    ar := binary.BigEndian.Uint16(buf[10:])

    return &DNSHeader {
        id, pack, qd, an, ns, ar,
    }
}

func (h *DNSHeader) IsQuery() bool {
    return h.Packed & 0b1000_0000_0000_0000 == 0
}

func main() {
    l, err := net.ListenPacket("udp", "127.0.0.1:4444")
    if err != nil {
        panic(err)
    }

    defer l.Close()

    for {
        buf := make([]byte, 16*6) // DNS header
        if err != nil {
            panic(err)
        }

        n, addr, err := l.ReadFrom(buf)
        if err != nil {
            panic(err)
        }

        go func(pc net.PacketConn, addr net.Addr, buf []byte) {
            header := NewDNSHeaderFromBytes(buf)

            fmt.Println(">>>> ", addr, " >> ", buf)
            fmt.Println(">>>> is query: ", header.IsQuery(), " QD count: ", header.QDcount, " AN: ", header.ANcount)
        }(l, addr, buf[:n])
    }
}
