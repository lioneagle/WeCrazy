package rtp

type RtpProfile struct {
	Used         bool
	PayloadType  byte
	Name         string
	MediaType    string
	HasClockRate bool
	ClockRate    uint32
	HasChannels  bool
	Channels     byte
}

// static rtp profiles from RFC3551
var StaticRtpProfiles = [95]RtpProfile{
	// static audio codecs
	0:  {Used: true, Name: "PCMU", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	3:  {Used: true, Name: "GSM", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	4:  {Used: true, Name: "G723", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	5:  {Used: true, Name: "DVI4", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	6:  {Used: true, Name: "DVI4", MediaType: "A", HasClockRate: true, ClockRate: 16000, HasChannels: true, Channels: 1},
	7:  {Used: true, Name: "LPC", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	8:  {Used: true, Name: "PCMA", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	9:  {Used: true, Name: "G722", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	10: {Used: true, Name: "L16", MediaType: "A", HasClockRate: true, ClockRate: 44100, HasChannels: true, Channels: 2},
	11: {Used: true, Name: "L16", MediaType: "A", HasClockRate: true, ClockRate: 44100, HasChannels: true, Channels: 1},
	12: {Used: true, Name: "QCELP", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	13: {Used: true, Name: "CN", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	14: {Used: true, Name: "MPA", MediaType: "A", HasClockRate: true, ClockRate: 90000, HasChannels: false},
	15: {Used: true, Name: "G728", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},
	16: {Used: true, Name: "DVI4", MediaType: "A", HasClockRate: true, ClockRate: 11025, HasChannels: true, Channels: 1},
	17: {Used: true, Name: "DVI4", MediaType: "A", HasClockRate: true, ClockRate: 22050, HasChannels: true, Channels: 1},
	18: {Used: true, Name: "G729", MediaType: "A", HasClockRate: true, ClockRate: 8000, HasChannels: true, Channels: 1},

	// static video codecs
	25: {Used: true, Name: "CelB", MediaType: "V", HasClockRate: true, ClockRate: 90000, HasChannels: false},
	26: {Used: true, Name: "JPEG", MediaType: "V", HasClockRate: true, ClockRate: 90000, HasChannels: false},
	28: {Used: true, Name: "nv", MediaType: "V", HasClockRate: true, ClockRate: 90000, HasChannels: false},
	31: {Used: true, Name: "H261", MediaType: "V", HasClockRate: true, ClockRate: 90000, HasChannels: false},
	32: {Used: true, Name: "MPV", MediaType: "V", HasClockRate: true, ClockRate: 90000, HasChannels: false},
	33: {Used: true, Name: "MP2T", MediaType: "AV", HasClockRate: true, ClockRate: 90000, HasChannels: false},
	34: {Used: true, Name: "H263", MediaType: "V", HasClockRate: true, ClockRate: 90000, HasChannels: false},
}

func GetStaticPayloadTypeName(payloadType byte) string {
	if payloadType > byte(len(StaticRtpProfiles)) {
		return "dynamic"
	}
	if !StaticRtpProfiles[payloadType].Used {
		return "unknown"
	}
	return StaticRtpProfiles[payloadType].Name
}
