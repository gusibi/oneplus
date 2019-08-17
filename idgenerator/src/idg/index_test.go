//
package idg

import (
	"testing"
)

func TestBinarySearch(t *testing.T) {
	values := []string{"a", "b", "c", "e", "f", "g"}
	t.Log(BinarySearch(values, "b"))
	t.Log(BinarySearch(values, "d"))
	t.Log(BinarySearch(values, "e"))
}
