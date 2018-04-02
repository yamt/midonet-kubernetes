package midonet

import (
	"net"
	"reflect"
	"testing"
)

func TestMac(t *testing.T) {
	actual := macForKey("hey")
	expected, err := net.ParseMAC("ac:ca:ba:fa:69:0b")
	if err != nil {
		t.Errorf("ParseMAC")
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
