package compute2

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

type Computer struct {
}

func (t *Computer) NetworkCompute(a, b int) (int, error) {
	// do something in remote computer
	c := a + b

	return c, nil
}

func (t *Computer) Compute(a, b int) (int, error) {
	sum, err := t.NetworkCompute(a, b)
	return sum, err
}

func TestCompute(t *testing.T) {
	var c *Computer
	gomonkey.ApplyMethod(reflect.TypeOf(c), "NetworkCompute", func(_ *Computer, a, b int) (int, error) {
		return 2999, nil
	})

	//defer patches.Reset()

	cp := &Computer{}
	sum, err := cp.Compute(1, 1)
	if sum != 2 || err != nil {
		t.Errorf("expected %v, got %v", 2, sum)
	}

}
