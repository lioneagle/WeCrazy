package rtp

import (
	"fmt"
	"os"
	"testing"

	"github.com/lioneagle/goutil/src/buffer"
	"github.com/lioneagle/goutil/src/test"
)

func TestRtpPacket(t *testing.T) {
	type result struct {
		version          byte
		padding          bool
		extension        bool
		csrcCount        byte
		marker           bool
		payloadType      byte
		sequence         uint16
		timestamp        uint32
		ssrc             uint32
		csrc             []uint32
		extensionProfile uint16
		extensions       []byte
		setExtensionsOk  bool
	}

	inputs := []struct {
		src *result
		ret *result
	}{
		{&result{version: 1, padding: true, extension: true, marker: true, payloadType: 8, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11, 12, 13, 14}, setExtensionsOk: true},
			&result{version: 1, padding: true, extension: true, marker: true, payloadType: 8, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11, 12, 13, 14}, setExtensionsOk: true}},
		{&result{version: 2, padding: false, extension: true, marker: false, payloadType: 0, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11, 12, 13}, setExtensionsOk: true},
			&result{version: 2, padding: false, extension: true, marker: false, payloadType: 0, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11, 12, 13, 0}, setExtensionsOk: true}},
		{&result{version: 2, padding: false, extension: true, marker: false, payloadType: 0, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11, 12}, setExtensionsOk: true},
			&result{version: 2, padding: false, extension: true, marker: false, payloadType: 0, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11, 12, 0, 0}, setExtensionsOk: true}},
		{&result{version: 2, padding: false, extension: true, marker: false, payloadType: 0, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11}, setExtensionsOk: true},
			&result{version: 2, padding: false, extension: true, marker: false, payloadType: 0, sequence: 2, timestamp: 100, csrcCount: 4, ssrc: 123456, csrc: []uint32{1, 2, 3, 4}, extensionProfile: 5678, extensions: []byte{11, 0, 0, 0}, setExtensionsOk: true}},
	}

	garbage := make([]byte, 1000)
	for i := 0; i < len(garbage); i++ {
		garbage[i] = 0xff
	}

	for i, v := range inputs {
		v := v
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			rtp := NewRtpPacket()
			rtp.CopyFromBytes(garbage)

			//rtp := NewRtpPacket(1000)
			ret := &result{}

			rtp.SetVersion(v.src.version)
			if v.src.padding {
				rtp.SetPadding()
			} else {
				rtp.ClearPadding()
			}
			if v.src.extension {
				rtp.SetExtensionBit()
			} else {
				rtp.ClearExtensionBit()
			}
			if v.src.marker {
				rtp.SetMarker()
			} else {
				rtp.ClearMarker()
			}
			rtp.SetCsrcCount(v.src.csrcCount)
			rtp.SetPayloadType(v.src.payloadType)
			rtp.SetSequence(v.src.sequence)
			rtp.SetTimestamp(v.src.timestamp)
			rtp.SetSsrc(v.src.ssrc)
			rtp.SetCsrc(v.src.csrc)
			if len(v.src.extensions) > 0 {
				ret.setExtensionsOk = rtp.SetExtension(v.src.extensionProfile, v.src.extensions)
			}

			ret.version = rtp.GetVersion()
			ret.padding = rtp.GetPadding() == 1
			ret.extension = rtp.GetExtensionBit() == 1
			ret.csrcCount = rtp.GetCsrcCount()
			ret.marker = rtp.GetMarker() == 1
			ret.payloadType = rtp.GetPayloadType()
			ret.sequence = rtp.GetSequence()
			ret.timestamp = rtp.GetTimestamp()
			ret.ssrc = rtp.GetSsrc()
			ret.csrc = rtp.GetCsrc()
			ret.extensionProfile = rtp.GetExtensionProfile()
			ret.extensions = rtp.GetExtension()

			fmt.Println(v.ret)

			ok, msg := test.DiffEx("", ret, v.ret)
			if !ok {
				rtp.Print(os.Stdout)
				t.Errorf("\n" + msg)
			}
		})
	}
}

func TestRtpPacketPrint(t *testing.T) {
	buf := buffer.NewByteBuffer(nil)

	rtp := NewRtpPacket()
	size := rtp.CalcLen(2, 2, 20)
	rtp.Alloc(size)

	rtp.SetVersion(3)
	rtp.SetPadding()
	rtp.SetExtensionBit()
	rtp.SetMarker()
	rtp.SetCsrcCount(2)
	rtp.SetPayloadType(8)
	rtp.SetSequence(1234)
	rtp.SetTimestamp(5678)
	rtp.SetSsrc(12345678)
	rtp.SetCsrc([]uint32{12, 13})
	rtp.SetExtension(4567, []byte{1, 2, 3, 4, 5, 6, 7})

	rtp.Print(buf)

	fmt.Println(buf)

	wanted := `11.. .... = version: 3
..1. .... = Padding: true
...1 .... = Extension: true
.... 0010 = CSRC count: 2
.1.. .... = Marker: true
Payload type: PCMA (8)
Sequence number: 1234
Timestamp: 5678
SSRC: 0x00bc614e (12345678)
CSRC:
[0]: 0x0000000c (12)
[1]: 0x0000000d (13)
Extension:
profile:0x11d7 (4567)
00000000h: 01 02 03 04 05 06 07 00                          ; ........
Payload:
00000000h: 07 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ; ................
00000010h: 00 00 00 00                                      ; ....
`
	test.EXPECT_EQ(t, buf.String(), wanted, "")

}

/*
func Test2(t *testing.T) {
	type result struct {
		y1 int
		Y2 error `break`
		y3 int
	}

	type testdata struct {
		x1  int
		x2  int
		Ret *result
	}

	inputs := []testdata{
		{11, 21, &result{31, errors.New("123"), 61}},
		{12, 22, &result{41, &MyError{10, 11}, 28}},
		{13, 23, &result{51, nil, 32}},
		{13, 23, &result{51, nil, 79}},
		{13, 23, &result{51, nil, 1}},
	}

	ok, msg := test.TestGroup(t, inputs, func(val interface{}) interface{} {
		input := val.(testdata)
		ret := result{}
		ret.y1 = input.x1 + 1

		return &ret
	})

	if !ok {
		t.Error("\n" + msg)
	}
}*/
