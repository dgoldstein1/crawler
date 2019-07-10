package crawler

import (
  "reflect"
  "testing"
)

// adopted taken from https://gist.github.com/samalba/6059502
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received '%v' (type %v), expected '%v' (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func AssertErrorEqual(t *testing.T, a error, b error) {
	if (a == nil || b == nil) {
		AssertEqual(t, a, b)
		return
	}
	if (a.Error() == b.Error()) {
		return
	}
	t.Errorf("Received '%v' (type %v), expected '%v' (type %v)", a.Error(), reflect.TypeOf(a), b.Error(), reflect.TypeOf(b))
}
