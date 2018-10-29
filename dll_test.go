package main

import "testing"

func TestDLL(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected int
	}{
		{
			name:     "none",
			expected: 0,
			source: `
				package main

				import "fmt"

				func main() {
					fmt.Println("Hello!")
				}
			`,
		},
		{
			name:     "for",
			expected: 1,
			source: `
				package main

				import "fmt"

				func main() {
					for i := 0; i < 5; i++ {
						defer fmt.Println("defer")
					}
				}
			`,
		},
		{
			name:     "range",
			expected: 1,
			source: `
				package main

				import "fmt"

				func main() {
					list := []int{1, 2, 3, 4, 5, 6, 7}

					for _, x := range list {
						defer fmt.Println(x)
					}
				}
			`,
		},
		{
			name:     "nested",
			expected: 1,
			source: `
				package main

				import "fmt"

				func main() {
					list := []int{1, 2, 3, 4, 5, 6, 7}

					for _, i := range list {
						for j := 0; j < i; j++ {
							defer fmt.Println(j)
						}
					}
				}
			`,
		},
		{
			name:     "block",
			expected: 1,
			source: `
				package main

				import "fmt"

				func main() {
					for i := 0; i < 5; i++ {
						{
							defer fmt.Println("defer")
						}
					}
				}
			`,
		},
		{
			name:     "if",
			expected: 1,
			source: `
				package main

				import "fmt"

				func main() {
					for i := 0; i < 5; i++ {
						if true {
							defer fmt.Println("defer")
						}
					}
				}
			`,
		},
		{
			name:     "funclit",
			expected: 0,
			source: `
				package main

				import "fmt"

				func main() {
					for i := 0; i < 5; i++ {
						func() {
							defer fmt.Println("defer")
						}()
					}
				}
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reports, err := gather(tt.source, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(reports) != tt.expected {
				t.Fatalf("expected %d reports, got %d", tt.expected, len(reports))
			}
		})
	}
}

func TestErrorParsing(t *testing.T) {
	source := `
	package main

	func main() {
		s = "missing quote
	}
	`

	_, err := gather(source, false)
	if err == nil {
		t.Error("expected error but got nil")
	}
}
