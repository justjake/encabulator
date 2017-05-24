package assert

import "reflect"
import "testing"

// Equal asserts that the two values given are equal. A custom formatting
// string may be supplied.
func Equal(t *testing.T, a, b interface{}, message ...string) {
	if reflect.DeepEqual(a, b) {
		return
	}

	if len(message) == 0 {
		message = []string{"%#v != %#v"}
	}

	t.Errorf(message[0], a, b)
}
