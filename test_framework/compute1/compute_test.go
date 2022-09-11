package compute1

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

func networkCompute(a, b int) (int, error) {
	// do something in remote computer
	c := a + b

	fmt.Println(1111111111111)
	return c, nil
}

func Compute(a, b int) (int, error) {
	sum, err := networkCompute(a, b)
	return sum, err
}

func TestCompute(t *testing.T) {
	patches := gomonkey.ApplyFunc(networkCompute, func(a, b int) (int, error) {
		return 20, nil
	})
	defer patches.Reset()
	sum, err := Compute(1, 1)
	fmt.Println(sum)
	if sum != 2 || err != nil {
		t.Errorf("expected %v, got %v", 2, sum)
	}
}
