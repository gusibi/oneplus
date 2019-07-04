package handler

import "testing"

func TestGetMobileAttribution(t *testing.T) {
	t.Log(GetMobileAttribution("18511451103"))
	t.Log(GetMobileAttribution("+8618511451103"))
}
