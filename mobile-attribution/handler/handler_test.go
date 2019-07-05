package handler

import (
	"fmt"
	"testing"
)

func TestGetMobileAttribution(t *testing.T) {
	// t.Log(GetMobileAttribution("18543741111"))
	// t.Log(GetMobileAttribution("+8618511111234"))
	t.Log(GetMobileAttribution("19543741111"))
	t.Log(GetMobileAttribution("19800001111"))
}

func BenchmarkGetMobileAttribution(t *testing.B) {
	for n := 0; n < t.N; n++ {
		mobile := fmt.Sprintf("185%04d1111", n%10000)
		GetMobileAttribution(mobile)
	}
}
