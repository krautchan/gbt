// brainfuck_test.go
package brainfuck

import (
	"testing"
)

func TestBrainfuckA(t *testing.T) {
	bf := NewBrainfuckInterpreter("++++++++++.", "")

	ret, err := bf.Start()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	if len(ret) != 1 {
		t.Logf("Expected size 1: %v", ret)
		t.Fail()
	}

	if ret[0] != 10 {
		t.Logf("Epected ret[0] 10: %v", ret)
		t.Fail()
	}
}

func TestBrainfuckLoop(t *testing.T) {
	bf := NewBrainfuckInterpreter("+++++[>++++++++++.<-]", "")
	ret, err := bf.Start()
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}

	for i, v := range []byte{10, 20, 30, 40, 50} {
		if ret[i] != v {
			t.Logf("Expected: [10, 20, 30, 40, 50]: %v", ret)
			t.Fail()
		}
	}
}

func TestBrainfuckHelloWorld(t *testing.T) {
	bf := NewBrainfuckInterpreter("++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.", "")
	ret, err := bf.Start()
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}

	for i, v := range []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 10} {
		if ret[i] != v {
			t.Logf("Expected: [72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 10]: %v", ret)
			t.Fail()
		}
	}
}

func TestBrainfuckInput(t *testing.T) {
	bf := NewBrainfuckInterpreter(",.", "a")
	ret, err := bf.Start()
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}

	if ret[0] != 97 {
		t.Logf("%v", ret)
		t.Fail()
	}
}

func TestBrainfuckHelloWorld2(t *testing.T) {
	bf := NewBrainfuckInterpreter(">+++++++++[<++++++++>-]<.>+++++++[<++++>-]<+.+++++++..+++.[-]>++++++++[<++++>-]<.#>+++++++++++[<+++++>-]<.>++++++++[<+++>-]<.+++.------.--------.[-]>++++++++[<++++>-]<+.[-]++++++++++.", "")

	ret, err := bf.Start()
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}

	for i, v := range []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 10} {
		if ret[i] != v {
			t.Logf("Expected: [72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 10]: %v", ret)
			t.Fail()
		}
	}
}

func TestBrainfuckFibonacci(t *testing.T) {
	bf := NewBrainfuckInterpreter("+++++++++++>+>>>>++++++++++++++++++++++++++++++++++++++++++++>++++++++++++++++++++++++++++++++<<<<<<[>[>>>>>>+>+<<<<<<<-]>>>>>>>[<<<<<<<+>>>>>>>-]<[>++++++++++[-<-[>>+>+<<<-]>>>[<<<+>>>-]+<[>[-]<[-]]>[<<[>>>+<<<-]>>[-]]<<]>>>[>>+>+<<<-]>>>[<<<+>>>-]+<[>[-]<[-]]>[<<+>>[-]]<<<<<<<]>>>>>[++++++++++++++++++++++++++++++++++++++++++++++++.[-]]++++++++++<[->-<]>++++++++++++++++++++++++++++++++++++++++++++++++.[-]<<<<<<<<<<<<[>>>+>+<<<<-]>>>>[<<<<+>>>>-]<-[>>.>.<<<[-]]<<[>>+>+<<<-]>>>[<<<+>>>-]<<[<+>-]>[<+>-]<<<-]", "")

	ret, err := bf.Start()
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}

	for i, v := range []byte{49, 44, 32, 49, 44, 32, 50, 44, 32, 51, 44, 32, 53, 44, 32, 56, 44, 32, 49, 51, 44, 32, 50, 49, 44, 32, 51, 52, 44, 32, 53, 53, 44, 32, 56, 57} {
		if ret[i] != v {
			t.Logf("Expected: [49, 44, 32, 49, 44, 32, 50, 44, 32, 51, 44, 32, 53, 44, 32, 56, 44, 32, 49, 51, 44, 32, 50, 49, 44, 32, 51, 52, 44, 32, 53, 53, 44, 32, 56, 57]: %v", string(ret))
			t.Fail()
		}
	}
}

func TestBrainfuckEndlessLoop(t *testing.T) {
	bf := NewBrainfuckInterpreter("+[]", "")

	_, err := bf.Start()
	if err == nil {
		t.Logf("Error expected")
		t.FailNow()
	}

	if err.Error() != "Too many loops" {
		t.Logf("Expected: To many loops error")
		t.Fail()
	}
}
