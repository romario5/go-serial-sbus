package sbus

import (
	"errors"
	"fmt"
	"io"

	serial "github.com/romario5/go-serial-common"
)

const (
	HEADER_BYTE     = 0x0F
	FOOTER_BYTE     = 0x00
	FOOTER2_BYTE    = 0x04
	CH17_MASK       = 0x01
	CH18_MASK       = 0x02
	LOST_FRAME_MASK = 0x04
	FAILSAFE_MASK   = 0x08
	CHANNELS_COUNT  = 18
	PACKET_LENGTH   = 25
	PAYLOAD_LENGTH  = 23
	HEADER_LENGTH   = 1
	FOOTER_LENGTH   = 1
)

type SBUS struct{}

func (s *SBUS) ReadPacket(r io.Reader) (any, error) {
	b := make([]byte, 1)
	var curByte byte
	var prevByte byte
	var state uint8 = 0
	packet := &serial.ChannelsPacket{}
	buffer := make([]byte, PACKET_LENGTH, PACKET_LENGTH)
	for {
		_, err := r.Read(b)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}

		curByte = b[0]
		if state == 0 {
			if (curByte == HEADER_BYTE) && ((prevByte == FOOTER_BYTE) ||
				((prevByte & 0x0F) == FOOTER2_BYTE)) {
				buffer[state] = curByte
				state++
			} else {
				state = 0
			}
		} else if state < PAYLOAD_LENGTH+HEADER_LENGTH {
			buffer[state] = curByte
			state++
		} else if state < PAYLOAD_LENGTH+HEADER_LENGTH+FOOTER_LENGTH {
			state = 0
			prevByte = curByte
			if (curByte == FOOTER_BYTE) || ((curByte & 0x0F) == FOOTER2_BYTE) {
				packet.ParseChannels(buffer[1:])

				if buffer[23]&CH17_MASK > 0 {
					packet.Channels[16] = 2000
				} else {
					packet.Channels[16] = 1000
				}

				if buffer[23]&CH18_MASK > 0 {
					packet.Channels[17] = 2000
				} else {
					packet.Channels[17] = 1000
				}

				packet.IsFailsafe = buffer[23]&FAILSAFE_MASK > 0
				return packet, nil
			} else {
				return nil, errors.New("Invalid packet structure")
			}
		} else {
			state = 0
		}
		prevByte = curByte
	}
}

func (s *SBUS) WritePacket(w io.Writer, p any) error {
	packet, ok := p.(serial.ChannelsPacket)
	if ok {
		data := make([]byte, PACKET_LENGTH, PACKET_LENGTH)
		data[0] = HEADER_BYTE
		packet.PackChannels(data[1:])
		data[24] = FOOTER_BYTE
		_, err := w.Write(data)
		return err
	}
	return errors.New("Only ChannelsPacket is allowed to write using SBUS protocol")
}
