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

				func main() {
					println("Hello!")
				}
			`,
		},
		{
			name:     "for",
			expected: 1,
			source: `
				package main

				func main() {
					for i := 0; i < 5; i++ {
						defer println("defer")
					}
				}
			`,
		},
		{
			name:     "range",
			expected: 1,
			source: `
				package main

				func main() {
					list := []int{1, 2, 3, 4, 5, 6, 7}
					for _, x := range list {
						defer println(x)
					}
				}
			`,
		},
		{
			name:     "nested",
			expected: 1,
			source: `
				package main

				func main() {
					list := []int{1, 2, 3, 4, 5, 6, 7}
					for _, i := range list {
						for j := 0; j < i; j++ {
							defer println(j)
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

				func main() {
					for i := 0; i < 5; i++ {
						{
							defer println("defer")
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

				func main() {
					for i := 0; i < 5; i++ {
						if true {
							defer println("defer")
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

				func main() {
					for i := 0; i < 5; i++ {
						func() {
							defer println("defer")
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

func Test_splitArrayIntoParts(t *testing.T) {
	getStringArray := func(amount int) []string {
		strings := make([]string, 0, amount)
		for i := 0; i < amount; i++ {
			strings = append(strings, "foo")
		}
		return strings
	}

	tests := []struct {
		name          string
		strings       []string
		parts         int
		expectedParts int
	}{
		{
			name:          "should split the array with one string into one part",
			strings:       getStringArray(1),
			parts:         1,
			expectedParts: 1,
		},
		{
			name:          "should split the array with one string into zero part",
			strings:       getStringArray(1),
			parts:         0,
			expectedParts: 1,
		},
		{
			name:          "should split the array with two strings into one part",
			strings:       getStringArray(2),
			parts:         1,
			expectedParts: 1,
		},
		{
			name:          "should split the array with two strings into four part",
			strings:       getStringArray(2),
			parts:         4,
			expectedParts: 2,
		},
		{
			name:          "should split the array with one string into two part",
			strings:       getStringArray(1),
			parts:         2,
			expectedParts: 1,
		},
		{
			name:          "should split the array with four strings into two part",
			strings:       getStringArray(4),
			parts:         2,
			expectedParts: 2,
		},
		{
			name:          "should split the array with two strings into three part",
			strings:       getStringArray(2),
			parts:         3,
			expectedParts: 2,
		},
		{
			name:          "should split the array with ten strings into three part",
			strings:       getStringArray(10),
			parts:         3,
			expectedParts: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := splitArrayIntoParts(test.strings, test.parts)

			if len(got) != test.expectedParts {
				t.Fatalf("Expect to split the array into '%d' but got '%d'", test.expectedParts, len(got))
			}

			for _, files := range got {
				if len(files) < 1 {
					t.Fatalf("Expected to contain at least on string but got none")
				}
			}
		})
	}

	t.Run("should split the empty array into one part", func(t *testing.T) {
		strings := []string{}
		parts := 1

		got := len(splitArrayIntoParts(strings, parts))
		want := 1

		if got != want {
			t.Fatalf("Expected a length of %d but got %d", want, got)
		}
	})
}
