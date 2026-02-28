package validator

import (
    "testing"
)

func TestIsValidSubscriptionID(t *testing.T) {
    type Args struct {
        id string
    }

    type testCase struct {
        name string
        Args Args
        Want bool
    }

    runTests := func(t *testing.T, tests []testCase) {
        t.Helper()
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                if got := IsValidSubscriptionID(tt.Args.id); got != tt.Want {
                    t.Errorf("IsValidSubscriptionID() = %v, Want %v (input: %q)", got, tt.Want, tt.Args.id)
                }
            })
        }
    }

    t.Run("Valid", func(t *testing.T) {
        tests := []testCase{
            {"Numeric", Args{"123"}, true},
            {"AlphaNumeric", Args{"sub-abc-123"}, true},
            {"LongID", Args{"a-very-long-subscription-id-0123456789"}, true},
        }
        runTests(t, tests)
    })

    t.Run("Invalid_ControlChars", func(t *testing.T) {
        tests := []testCase{
            {"Empty", Args{""}, false},
            {"NullByte", Args{"abc\x00def"}, false},
            {"Ctrl1", Args{"abc\x01"}, false},
            {"DEL", Args{"abc\x7f"}, false},
            {"ContainsSpace", Args{"sub id"}, false},
        }
        runTests(t, tests)
    })

    t.Run("Invalid_URLChars", func(t *testing.T) {
        tests := []testCase{
            {"Slash", Args{"a/b"}, false},
            {"Question", Args{"a?b"}, false},
            {"Hash", Args{"a#b"}, false},
        }
        runTests(t, tests)
    })
}
