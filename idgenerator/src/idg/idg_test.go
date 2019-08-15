package idg

import "testing"

func TestIDNumberFill(t *testing.T) {
	// 100000192201010030
	t.Log(IDNumberFill("10000019220101003"))
	// 100000198001010033
	t.Log(IDNumberFill("10000019800101003"))
}

func TestId2Md5(t *testing.T) {
	// 062a6a230b49f2678999226f070e3e08
	t.Log(Md5("100000192201010030"))
	// 988feb91c2847d74b63a4a865691b763
	t.Log(Md5("100000198001010033"))
}
