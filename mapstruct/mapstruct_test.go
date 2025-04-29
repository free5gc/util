package mapstruct

import (
	"reflect"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

type TestProfile struct {
	RecoveryTime *time.Time `mapstructure:"RecoveryTime"`
}

func TestDecode(t *testing.T) {
	date, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Errorf("time parse fail: %+v", err)
	}

	testCases := []struct {
		description    string
		source         map[string]interface{}
		decodeError    bool
		resultStr      string
		expectedTarget TestProfile
	}{
		{
			description: "Param: RFC3339 time format string",
			source: map[string]interface{}{
				"RecoveryTime": date.Format(time.RFC3339),
			},
			decodeError: false,
			resultStr:   "Time field check should pass",
			expectedTarget: TestProfile{
				RecoveryTime: &date,
			},
		},
		{
			description: "Param: no time field entry",
			source:      map[string]interface{}{},
			decodeError: false,
			resultStr:   "Time field should be nil",
			expectedTarget: TestProfile{
				RecoveryTime: nil,
			},
		},
		{
			description: "Param: ANSIC time format string",
			source: map[string]interface{}{
				"RecoveryTime": date.Format(time.ANSIC),
			},
			decodeError: true,
			resultStr:   "Time field should be nil",
			expectedTarget: TestProfile{
				RecoveryTime: nil,
			},
		},
	}

	convey.Convey("Decode Test", t, func() {
		for _, tc := range testCases {
			convey.Convey(tc.description, func() {
				var target TestProfile
				err := Decode(tc.source, &target)
				if tc.decodeError {
					convey.Convey("Decode should fail", func() {
						convey.So(err, convey.ShouldNotBeNil)
					})
				} else {
					convey.Convey("Decode should succeed", func() {
						convey.So(err, convey.ShouldBeNil)
					})
				}
				convey.Convey(tc.resultStr, func() {
					convey.So(true, convey.ShouldEqual, reflect.DeepEqual(target.RecoveryTime, tc.expectedTarget.RecoveryTime))
				})
			})
		}
	})
}
