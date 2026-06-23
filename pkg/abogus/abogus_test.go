package abogus

import (
	"testing"
)

func TestGenerateABogus(t *testing.T) {
	params := "device_platform=webapp&aid=6383&channel=channel_pc_web&pc_client_type=1&version_code=190500&version_name=19.5.0&cookie_enabled=true&browser_language=zh-CN&browser_platform=Win32&browser_name=Firefox&browser_online=true&engine_name=Gecko&os_name=Windows&os_version=10&platform=PC&screen_width=1920&screen_height=1080&browser_version=124.0&engine_version=122.0.0.0&cpu_core_num=12&device_memory=8&aweme_id=7345492945006595379"
	method := "GET"
	startTime := int64(1718200000000)
	endTime := int64(1718200000005)
	rn1 := 1234.56
	rn2 := 5678.90
	rn3 := 9012.34

	expected := "E7mhBdugDifihdWk5RxLfY3q6VWVYmBD0SVkMD2fn-DOOg39HMYh9exooCivRm8jNs/DIeEjy4hbT3ohrQ2y0Hwf9W0L/25ksDSkKl5Q5xSSs1X9eghgJ04qmkt5SMx2RvB-rOXmqhZHKRbp09oHmhK4b1dzFgf3qJLzbj=="

	result := GenerateABogusWithOptions(params, method, startTime, endTime, rn1, rn2, rn3)
	if result != expected {
		t.Errorf("ABogus mismatch:\nExpected: %s\nActual:   %s", expected, result)
	} else {
		t.Logf("Success! Generated exact match: %s", result)
	}
}
