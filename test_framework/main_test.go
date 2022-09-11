package main_test

import (
	"encoding/json"
	"fmt"

	"testing"
	"unicode"
)

func TestIs(t *testing.T) {
	fmt.Println(123)
}

const (
	Left = iota
	Top
	Right
	Bottom
)

func move(cmd string, x0, y0, z0 int) (x, y, z int) {
	x, y, z = x0, y0, z0
	repeat := 0
	repeatCmd := ""
	for _, s := range cmd {
		switch {
		case unicode.IsNumber(s):
			repeat = repeat*10 + (int(s) - '0')
		case s == ')':
			for i := 0; i < repeat; i++ {
				x, y, z = move(repeatCmd, x, y, z)
			}
			repeat = 0
			repeatCmd = ""
		case repeat > 0 && s != '(' && s != ')':
			repeatCmd += string(s)
		case s == 'L':
			z = (z + 1) % 4
		case s == 'R':
			z = (z - 1 + 4) % 4
		case s == 'F':
			switch {
			case z == Left || z == Right:
				x = x + z - 1
			case z == Top || z == Bottom:
				y = y - z + 2
			}
		case s == 'B':
			switch {
			case z == Left || z == Right:
				x = x - z - 1
			case z == Top || z == Bottom:
				y = y + z - 2
			}
		}
	}
	return
}

func TestMove(t *testing.T) {
	fmt.Println(move("R2(LF)", 0, 0, Top))
}

type People struct {
	Name string `json:"name"`
}

func TestPeople(t *testing.T) {
	js := `{
		"name":"11"
	}`
	var p People
	err := json.Unmarshal([]byte(js), &p)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("people: ", p)
	fmt.Println(p.Name)
}

type Student struct {
	name string
}

func TestTain(t *testing.T) {
	m := map[string]Student{"people": {"zhoujielun"}}
	fmt.Println(m["people"].name)
	//p.name = "ddd"
	//m["people"].name = "wuyanzu"
}
