package main

import "testing"

func TestErrorParsing(t *testing.T) {
	source := `package main

	func main() {
		s = "missing quote
	}`

	_, err := gather(source, false)
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestNoDefer(t *testing.T) {
	source := `package main

	import "fmt"

	func main() {
		fmt.Println("Hello!")
	}`

	reports, err := gather(source, false)
	if err != nil {
		t.Errorf("expected no error; got %s", err)
	}

	if len(reports) > 0 {
		t.Errorf("expected no reports; got %d", len(reports))
	}
}

func TestSingleForLoopWithDefer(t *testing.T) {
	source := `package main

	import "fmt"

	func main() {
		for i := 0; i < 5; i++ {
			defer fmt.Println("defer")
		}
	}`

	reports, err := gather(source, false)
	if err != nil {
		t.Errorf("expected no error; got %s", err)
	}

	if len(reports) != 1 {
		t.Errorf("expected one report; got %d", len(reports))
	}
}

func TestNestedDefer(t *testing.T) {
	source := `package main

	import "fmt"

	func main() {
		list := []int{1, 2, 3, 4, 5, 6, 7}

		for _, i := range list {
			for j := 0; j < i; j++ {
				defer fmt.Println(j)
			}
		}
	}`

	reports, err := gather(source, false)
	if err != nil {
		t.Errorf("expected no error; got %s", err)
	}

	if len(reports) != 1 {
		t.Errorf("expected one report; got %d", len(reports))
	}
}

func TestBlockDefer(t *testing.T) {
	source := `package main

	import "fmt"

	func main() {
		for i := 0; i < len(list); i++ {
			{
				defer fmt.Println(i)
			}
		}
	}`

	reports, err := gather(source, false)
	if err != nil {
		t.Errorf("expected no error; got %s", err)
	}

	if len(reports) != 1 {
		t.Errorf("expected one report; got %d", len(reports))
	}

}

func TestIfDefer(t *testing.T) {
	source := `package main

	import "fmt"

	func main() {
		for i := 0; i < len(list); i++ {
			if true {
				defer fmt.Println(i)
			}
		}
	}`

	reports, err := gather(source, false)
	if err != nil {
		t.Errorf("expected no error; got %s", err)
	}

	if len(reports) != 1 {
		t.Errorf("expected one report; got %d", len(reports))
	}

}

func TestRangeDefer(t *testing.T) {
	source := `package main

	import "fmt"

	func main() {

		list := []int{1, 2, 3}

		for _, x := range list {
			defer fmt.Println(x)
		}
	}`

	reports, err := gather(source, false)
	if err != nil {
		t.Errorf("expected no error; got %s", err)
	}

	if len(reports) != 1 {
		t.Errorf("expected one report; got %d", len(reports))
	}

}
