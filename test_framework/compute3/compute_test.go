package compute3

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

func compute(a, b int) (int, error) {
	return a + b, nil
}

func TestFunc(t *testing.T) {
	info1 := "2"
	info2 := "3"
	info3 := "4"
	outputs := []gomonkey.OutputCell{
		{Values: gomonkey.Params{info1, nil}}, // 模拟函数的第1次输出
		{Values: gomonkey.Params{info2, nil}}, // 模拟函数的第2次输出
		{Values: gomonkey.Params{info3, nil}}, // 模拟函数的第3次输出
	}
	patches := gomonkey.ApplyFuncSeq(compute, outputs)
	defer patches.Reset()

	output, err := compute(1, 1)
	if output != 2 || err != nil {
		t.Errorf("expected %v, got %v", 2, output)
	}

	output, err = compute(1, 2)
	if output != 3 || err != nil {
		t.Errorf("expected %v, got %v", 2, output)
	}

	output, err = compute(1, 3)
	if output != 4 || err != nil {
		t.Errorf("expected %v, got %v", 2, output)
	}

}
