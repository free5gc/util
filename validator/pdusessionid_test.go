package validator

import (
	"testing"
)

func TestIsValidPduSessionID(t *testing.T) {
	type args struct {
		pduSessionID string
	}
	type testCase struct {
		name string
		args args
		want bool
	}

	// Helper function similar to your reference
	runTests := func(t *testing.T, tests []testCase) {
		t.Helper()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IsValidPduSessionID(tt.args.pduSessionID); got != tt.want {
					t.Errorf("IsValidPduSessionID() = %v, want %v (input: %q)", got, tt.want, tt.args.pduSessionID)
				}
			})
		}
	}

	// 1. Valid Range Tests (TS 29.571: integer 0 to 255)
	t.Run("Check_Valid_Range", func(t *testing.T) {
		tests := []testCase{
			{"Valid Min Value (0)", args{"0"}, true},
			{"Valid Max Value (255)", args{"255"}, true},
			{"Valid Mid Value", args{"128"}, true},
			{"Valid Single Digit", args{"5"}, true},
			// strconv.Atoi handles positive sign usually, technically valid integer
			{"Valid with Plus Sign", args{"+10"}, true},
		}
		runTests(t, tests)
	})

	// 2. Out of Range Tests (Boundary Analysis)
	t.Run("Check_Out_Of_Range", func(t *testing.T) {
		tests := []testCase{
			{"Invalid Negative (-1)", args{"-1"}, false},
			{"Invalid Negative Large", args{"-100"}, false},
			{"Invalid Max + 1 (256)", args{"256"}, false},
			{"Invalid Large Value", args{"9999"}, false},
			// Integer overflow for 64/32 bit (Atoi returns error)
			{"Invalid Huge Number", args{"99999999999999999999"}, false},
		}
		runTests(t, tests)
	})

	// 3. Invalid Format & Security Tests
	t.Run("Check_Invalid_Format_And_Security", func(t *testing.T) {
		tests := []testCase{
			{"Empty String", args{""}, false},
			{"Non-digit (Alphabet)", args{"abc"}, false},
			{"Mixed Alphanumeric", args{"12a"}, false},
			// Atoi rejects floats
			{"Float format (Dot)", args{"12.5"}, false},
			// Atoi interprets as 0 then stops or errors depending on context, usually errors on 'x'
			{"Hex format (0xFF)", args{"0xFF"}, false},
			{"Whitespace only", args{"   "}, false},
			// Atoi rejects spaces
			{"Leading/Trailing Space", args{" 123 "}, false},

			// [Security] Null Byte Injection check
			{"Security: Null Byte Middle", args{"1\x002"}, false},
			{"Security: Null Byte End", args{"12\x00"}, false},
			{"Security: Raw Null Bytes", args{"\x00\x00"}, false},
		}
		runTests(t, tests)
	})
}
