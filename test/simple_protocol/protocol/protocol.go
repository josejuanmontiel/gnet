package test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/panjf2000/gnet/v2"
)

func HandleMessage(c net.Conn) {
	s, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("HandleMessage could not read: %w", err))
	}
	s = strings.Trim(s, "\n")
	c.Write([]byte(s + " whatever\n"))
}

var ErrIncompletePacket = errors.New("incomplete packet")

const (
	magicNumber     = 1314
	magicNumberSize = 2
	bodySize        = 4
)

var magicNumberBytes []byte

func init() {
	magicNumberBytes = make([]byte, magicNumberSize)
	binary.BigEndian.PutUint16(magicNumberBytes, uint16(magicNumber))
}

type SimpleCodec struct{}

func (codec SimpleCodec) Encode(buf []byte) ([]byte, error) {
	bodyOffset := magicNumberSize + bodySize
	msgLen := bodyOffset + len(buf)

	data := make([]byte, msgLen)
	copy(data, magicNumberBytes)

	binary.BigEndian.PutUint32(data[magicNumberSize:bodyOffset], uint32(len(buf)))
	copy(data[bodyOffset:msgLen], buf)
	return data, nil
}

/*
func (codec *SimpleCodec) Decode(c gnet.Conn) ([]byte, error) {
	bodyOffset := magicNumberSize + bodySize
	buf, _ := c.Peek(bodyOffset)
	if len(buf) < bodyOffset {
		return nil, ErrIncompletePacket
	}

	if !bytes.Equal(magicNumberBytes, buf[:magicNumberSize]) {
		return nil, errors.New("invalid magic number")
	}

	bodyLen := binary.BigEndian.Uint32(buf[magicNumberSize:bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)

	return buf[bodyOffset:msgLen], nil
}
*/
func (codec *SimpleCodec) Decode(c gnet.Conn) ([]byte, error) {
	bodyOffset := magicNumberSize + bodySize

	var buf []byte
	var msgLen int
	v, ok := c.(gnet.GnetReader)
	if ok {
		buf, _ = v.Peek(bodyOffset)
		if len(buf) < bodyOffset {
			return nil, ErrIncompletePacket
		}
		if !bytes.Equal(magicNumberBytes, buf[:magicNumberSize]) {
			return nil, errors.New("invalid magic number")
		}

		bodyLen := binary.BigEndian.Uint32(buf[magicNumberSize:bodyOffset])
		msgLen = bodyOffset + int(bodyLen)
		if v.InboundBuffered() < msgLen {
			return nil, ErrIncompletePacket
		}
		buf, _ = v.Peek(msgLen)
		_, _ = v.Discard(msgLen)

		return buf[bodyOffset:msgLen], nil

	} else {
		// WIP TODO ...
		return buf, nil
	}

}

func (codec SimpleCodec) Unpack(buf []byte) ([]byte, error) {
	bodyOffset := magicNumberSize + bodySize
	if len(buf) < bodyOffset {
		return nil, ErrIncompletePacket
	}

	if !bytes.Equal(magicNumberBytes, buf[:magicNumberSize]) {
		return nil, errors.New("invalid magic number")
	}

	bodyLen := binary.BigEndian.Uint32(buf[magicNumberSize:bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	if len(buf) < msgLen {
		return nil, ErrIncompletePacket
	}

	return buf[bodyOffset:msgLen], nil
}
