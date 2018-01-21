package rtp

func init() {
	for i := 0; i < len(StaticRtpProfiles); i++ {
		StaticRtpProfiles[i].PayloadType = byte(i)
	}
}
