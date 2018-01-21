package rtp

import (
	"encoding/binary"
	"fmt"
	"io"
	//"os"

	"github.com/lioneagle/goutil/src/buffer"
)

const (
	RTP_VERSION_MARSK      = 0xC0
	RTP_PADDING_MARSK      = 0x20
	RTP_EXTENSION_MARSK    = 0x10
	RTP_CSRC_COUNT_MARSK   = 0x0F
	RTP_MARKER_MARSK       = 0x80
	RTP_PAYLOAD_TYPE_MARSK = 0x7F

	RTP_HEADER_LEN = 12
)

const (
	RTP_MARKER_OFFSET       = 1
	RTP_PAYLOAD_TYPE_OFFSET = 1
	RTP_SEQUENCE_OFFSET     = 2
	RTP_TIMESTAMP_OFFSET    = 4
	RTP_SSRC_OFFSET         = 8
)

type RtpPacket struct {
	data []byte
}

func NewRtpPacket() *RtpPacket {
	return &RtpPacket{}
}

func (this *RtpPacket) Alloc(size int) {
	this.data = make([]byte, size)
}

func (this *RtpPacket) Reset() {
	this.data = this.data[:0]
}

func (this *RtpPacket) CalcLen(csrcNum, extensionNum, payloadLength int) int {
	return this.CalcHeaderLen(csrcNum, extensionNum) + payloadLength
}

func (this *RtpPacket) CalcHeaderLen(csrcNum, extensionNum int) int {
	length := RTP_HEADER_LEN
	if csrcNum > 0 {
		length += csrcNum * 4
	}
	if extensionNum > 0 {
		length += 2 + extensionNum*4
	}
	return length
}

func (this *RtpPacket) HeaderLen() int {
	csrcNum := int(this.GetCsrcCount())
	extensionNum := 0
	if this.GetExtensionBit() == 1 {
		extensionNum = this.GetExtensionNum(csrcNum)
	}
	return this.CalcHeaderLen(csrcNum, extensionNum)
}

func (this *RtpPacket) PayloadLen() int {
	return len(this.data) - this.HeaderLen()
}

func (this *RtpPacket) GetPayload() []byte {
	return this.data[this.HeaderLen():]
}

func (this *RtpPacket) GetVersion() byte {
	return (this.data[0] & RTP_VERSION_MARSK) >> 6
}

func (this *RtpPacket) SetVersion(val byte) {
	this.data[0] &^= RTP_VERSION_MARSK
	this.data[0] |= (val << 6)
}

func (this *RtpPacket) GetPadding() byte {
	return (this.data[0] & RTP_PADDING_MARSK) >> 5
}

func (this *RtpPacket) SetPadding() {
	this.data[0] |= RTP_PADDING_MARSK
}

func (this *RtpPacket) ClearPadding() {
	this.data[0] &^= RTP_PADDING_MARSK
}

func (this *RtpPacket) GetExtensionBit() byte {
	return (this.data[0] & RTP_EXTENSION_MARSK) >> 4
}

func (this *RtpPacket) SetExtensionBit() {
	this.data[0] |= RTP_EXTENSION_MARSK
}

func (this *RtpPacket) ClearExtensionBit() {
	this.data[0] &^= RTP_EXTENSION_MARSK
}

func (this *RtpPacket) GetCsrcCount() byte {
	return this.data[0] & RTP_CSRC_COUNT_MARSK
}

func (this *RtpPacket) SetCsrcCount(val byte) {
	this.data[0] &^= RTP_CSRC_COUNT_MARSK
	this.data[0] |= (val & RTP_CSRC_COUNT_MARSK)
}

func (this *RtpPacket) GetMarker() byte {
	return (this.data[RTP_MARKER_OFFSET] & RTP_MARKER_MARSK) >> 7
}

func (this *RtpPacket) SetMarker() {
	this.data[RTP_MARKER_OFFSET] |= RTP_MARKER_MARSK
}

func (this *RtpPacket) ClearMarker() {
	this.data[RTP_MARKER_OFFSET] &^= RTP_MARKER_MARSK
}

func (this *RtpPacket) GetPayloadType() byte {
	return this.data[RTP_PAYLOAD_TYPE_OFFSET] & RTP_PAYLOAD_TYPE_MARSK
}

func (this *RtpPacket) SetPayloadType(val byte) {
	this.data[RTP_PAYLOAD_TYPE_OFFSET] &^= RTP_PAYLOAD_TYPE_MARSK
	this.data[RTP_PAYLOAD_TYPE_OFFSET] |= (val & RTP_PAYLOAD_TYPE_MARSK)
}

func (this *RtpPacket) GetSequence() uint16 {
	return binary.BigEndian.Uint16(this.data[RTP_SEQUENCE_OFFSET:])
}

func (this *RtpPacket) SetSequence(val uint16) {
	binary.BigEndian.PutUint16(this.data[RTP_SEQUENCE_OFFSET:], val)
}

func (this *RtpPacket) GetTimestamp() uint32 {
	return binary.BigEndian.Uint32(this.data[RTP_TIMESTAMP_OFFSET:])
}

func (this *RtpPacket) SetTimestamp(val uint32) {
	binary.BigEndian.PutUint32(this.data[RTP_TIMESTAMP_OFFSET:], val)
}

func (this *RtpPacket) GetSsrc() uint32 {
	return binary.BigEndian.Uint32(this.data[RTP_SSRC_OFFSET:])
}

func (this *RtpPacket) SetSsrc(val uint32) {
	binary.BigEndian.PutUint32(this.data[RTP_SSRC_OFFSET:], val)
}

func (this *RtpPacket) GetCsrc() (csrc []uint32) {
	num := int(this.GetCsrcCount())

	for i := 0; i < num; i++ {
		csrc = append(csrc, binary.BigEndian.Uint32(this.data[RTP_HEADER_LEN+i*4:]))
	}
	return csrc
}

func (this *RtpPacket) SetCsrc(csrc []uint32) {
	num := len(csrc)
	this.SetCsrcCount(byte(num))
	for i := 0; i < num; i++ {
		binary.BigEndian.PutUint32(this.data[RTP_HEADER_LEN+i*4:], csrc[i])
	}
}

func (this *RtpPacket) GetExtension() []byte {
	if this.GetExtensionBit() == 0 {
		return nil
	}

	offset := RTP_HEADER_LEN + int(this.GetCsrcCount())*4 + 2
	num := int(binary.BigEndian.Uint16(this.data[offset:])) * 4

	return this.data[offset+2 : offset+2+num]
}

func (this *RtpPacket) GetExtensionNum(csrcNum int) int {
	offset := RTP_HEADER_LEN + csrcNum*4 + 2
	return int(binary.BigEndian.Uint16(this.data[offset:]))
}

func (this *RtpPacket) GetExtensionProfile() uint16 {
	if this.GetExtensionBit() == 0 {
		return 0
	}
	offset := RTP_HEADER_LEN + int(this.GetCsrcCount())*4
	return binary.BigEndian.Uint16(this.data[offset:])
}

func (this *RtpPacket) SetExtension(profile uint16, val []byte) bool {
	offset := RTP_HEADER_LEN + int(this.GetCsrcCount())*4
	num := len(val)
	pad := num & 0x3
	if pad != 0 {
		pad = 4 - pad
	}

	if (offset + 4 + num + pad) > len(this.data) {
		return false
	}

	binary.BigEndian.PutUint16(this.data[offset:], profile)
	binary.BigEndian.PutUint16(this.data[offset+2:], uint16((num+pad)/4))

	copy(this.data[offset+4:], val)

	if pad == 0 {
		return true
	}

	paddOffset := offset + 4 + num
	if pad == 1 {
		this.data[paddOffset] = 0
	} else if pad == 2 {
		this.data[paddOffset] = 0
		this.data[paddOffset+1] = 0

	} else if pad == 3 {
		this.data[paddOffset] = 0
		this.data[paddOffset+1] = 0
		this.data[paddOffset+2] = 0
	}

	return true
}

func (this *RtpPacket) CopyFromBytes(data []byte) {
	this.Reset()
	this.data = append(this.data, data...)
}

func (this *RtpPacket) CopyToBytes(data []byte) {
	data = append(data, this.data...)
}

func (this *RtpPacket) CopyToByteBuffer(buf *buffer.ByteBuffer) {
	buf.Reset()
	buf.Write(this.data)
}

func (this *RtpPacket) Print(w io.Writer) {
	fmt.Fprintf(w, "%02b.. .... = version: %d\n", this.GetVersion(), this.GetVersion())
	fmt.Fprintf(w, "..%01b. .... = Padding: %v\n", this.GetPadding(), this.GetPadding() == 1)
	fmt.Fprintf(w, "...%01b .... = Extension: %v\n", this.GetExtensionBit(), this.GetExtensionBit() == 1)
	fmt.Fprintf(w, ".... %04b = CSRC count: %v\n", this.GetCsrcCount(), this.GetCsrcCount())
	fmt.Fprintf(w, ".%01b.. .... = Marker: %v\n", this.GetMarker(), this.GetMarker() == 1)
	fmt.Fprintf(w, "Payload type: %s (%v)\n", GetStaticPayloadTypeName(this.GetPayloadType()), this.GetPayloadType())
	fmt.Fprintf(w, "Sequence number: %v\n", this.GetSequence())
	fmt.Fprintf(w, "Timestamp: %v\n", this.GetTimestamp())
	fmt.Fprintf(w, "SSRC: 0x%08x (%d)\n", this.GetSsrc(), this.GetSsrc())
	if this.GetCsrcCount() > 0 {
		csrc := this.GetCsrc()
		fmt.Fprintf(w, "CSRC:\n")
		for i, v := range csrc {
			fmt.Fprintf(w, "[%d]: 0x%08x (%d)\n", i, v, v)
		}
	}
	if this.GetExtensionBit() == 1 {
		profile := this.GetExtensionProfile()
		extension := this.GetExtension()
		fmt.Fprintf(w, "Extension:\n")
		fmt.Fprintf(w, "profile:0x%04x (%d)\n", profile, profile)
		buffer.PrintAsHex(w, extension, 0, len(extension))
	}
	payload := this.GetPayload()
	if len(payload) > 0 {
		fmt.Fprintf(w, "Payload:\n")
		buffer.PrintAsHex(w, payload, 0, len(payload))
	}
}
