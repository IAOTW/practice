package compute5

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

var flag bool

func IsEnabled() bool {
	return flag
}

func Compute(a, b int) int {
	if IsEnabled() {
		return a + b
	}

	return a - b
}

func TestCompute(t *testing.T) {
	patches := gomonkey.ApplyFunc(IsEnabled, func() bool {
		return true
	})

	defer patches.Reset()

	sum := Compute(1, 1)
	if sum != 2 {
		t.Errorf("expected %v, got %v", 2, sum)
	}

}
