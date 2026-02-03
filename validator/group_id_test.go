package validator

import (
	"testing"
)

func TestValidateGroupIdFormat(t *testing.T) {
	type args struct {
		groupId string
	}
	type testCase struct {
		name    string
		args    args
		wantErr bool
	}

	runTests := func(t *testing.T, tests []testCase) {
		t.Helper()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := ValidateGroupIdFormat(tt.args.groupId)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateGroupIdFormat() error = %v, wantErr %v (input: %q)", err, tt.wantErr, tt.args.groupId)
				}
			})
		}
	}

	// Valid Group ID test
	t.Run("Check_Valid_GroupId", func(t *testing.T) {
		tests := []testCase{
			{"Valid GroupId (Standard)", args{"prefix-001-01-localid"}, false},
			{"Valid GroupId (MNC 3 digits)", args{"prefix-001-001-localid"}, false},
			{"Valid GroupId (Long LocalId)", args{"prefix-123-456-very-long-local-identifier"}, false},
			{"Valid GroupId (Numeric Prefix)", args{"123-456-78-localid"}, false},
			{"Valid GroupId (Alphanumeric)", args{"abc-123-45-xyz789"}, false},
		}
		runTests(t, tests)
	})

	// Invalid Group ID test - MCC validation
	t.Run("Check_Invalid_MCC", func(t *testing.T) {
		tests := []testCase{
			{"Invalid MCC (Too short)", args{"prefix-01-01-localid"}, true},
			{"Invalid MCC (Too long)", args{"prefix-0001-01-localid"}, true},
			{"Invalid MCC (Non-digits)", args{"prefix-abc-01-localid"}, true},
		}
		runTests(t, tests)
	})

	// Invalid Group ID test - MNC validation
	t.Run("Check_Invalid_MNC", func(t *testing.T) {
		tests := []testCase{
			{"Invalid MNC (Too short)", args{"prefix-001-1-localid"}, true},
			{"Invalid MNC (Too long)", args{"prefix-001-0001-localid"}, true},
			{"Invalid MNC (Non-digits)", args{"prefix-001-ab-localid"}, true},
		}
		runTests(t, tests)
	})

	// Invalid Group ID test - Structure validation
	t.Run("Check_Invalid_Structure", func(t *testing.T) {
		tests := []testCase{
			{"Missing Prefix", args{"-001-01-localid"}, true},
			{"Missing MCC", args{"prefix--01-localid"}, true},
			{"Missing MNC", args{"prefix-001--localid"}, true},
			{"Missing LocalId", args{"prefix-001-01-"}, true},
			{"Too Few Parts", args{"prefix-001-01"}, true},
			{"Only Two Parts", args{"prefix-001"}, true},
			{"Only One Part", args{"prefix"}, true},
			{"Empty String", args{""}, true},
		}
		runTests(t, tests)
	})

	// Invalid Group ID test - General invalid cases
	t.Run("Check_General_Invalid", func(t *testing.T) {
		tests := []testCase{
			{"No Dashes", args{"prefix001001localid"}, true},
			{"Extra Dash at Start", args{"-prefix-001-01-localid"}, true},
			{"Double Dash", args{"prefix--001-01-localid"}, true},
			{"Special Characters in MCC", args{"prefix-@01-01-localid"}, true},
			{"Special Characters in MNC", args{"prefix-001-@1-localid"}, true},
			{"Spaces", args{"prefix - 001 - 01 - localid"}, true},
		}
		runTests(t, tests)
	})
}
