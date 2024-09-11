package object

import "testing"

type HashKeyTest struct {
	expected   Hashable
	actual     Hashable
	unexpected Hashable
}

func TestHashKey(t *testing.T) {
	tests := []HashKeyTest{
		{
			&String{
				Value: "hello world",
			},
			&String{
				Value: "hello world",
			},
			&String{
				Value: "goodbye moon",
			},
		},
		{
			&Boolean{
				Value: true,
			},
			&Boolean{
				Value: true,
			},
			&Boolean{
				Value: false,
			},
		},
		{
			&Integer{
				Value: 1,
			},
			&Integer{
				Value: 1,
			},
			&Integer{
				Value: 2,
			},
		},
	}

	for i, test := range tests {
		if test.expected.HashKey() != test.actual.HashKey() {
			t.Errorf("test[%d] - expected: %q actual: %q", i, test.expected.HashKey(), test.actual.HashKey())
			continue
		}

		if test.expected.HashKey() == test.unexpected.HashKey() {
			t.Errorf("test[%d] - expected: not equal but was: %q", i, test.unexpected.HashKey())
			continue
		}
	}
}
