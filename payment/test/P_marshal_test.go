package payment_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/TunnelWork/Ulysses.Lib/payment"
)

func typeEqual(a, b interface{}) bool {
	return fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
}

func TestMarshalUnmarshalP(t *testing.T) {
	originalP := payment.P{
		"String":        "Val1",
		"Boolean":       true,
		"BooleanString": "false",
		"EmbeddedP": payment.P{
			"EmbeddedString":    "Val2",
			"EmbeddedIntString": "1234",
			"EmbeddedInt":       1234,
		},
	}

	Pstr := originalP.String()
	recoverP := payment.Pify(Pstr)

	if !reflect.DeepEqual(originalP, recoverP) {
		t.Errorf("They are not equal.")
	}
}
