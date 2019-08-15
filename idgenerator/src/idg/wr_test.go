package idg

import (
	"testing"
)

func TestMd52Uint64(t *testing.T) {
	t.Log(Md52Uint64("ffffffffffffffffffffffffffffffff"))
}

func TestUint642Md5(t *testing.T) {
	t.Log(Uint642Md5(18446744073709551615, 18446744073709551615))
}
