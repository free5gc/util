package nas

import (
	"testing"

	nasie "github.com/free5gc/nas/ie"
	"github.com/free5gc/nas/message"
)

// The 5GMM message types of TS 24.501 table 9.7.1 all sit below this value,
// the 5GSM ones start at 0xc1. github.com/free5gc/nas encodes the same split
// in its MsgType constants but exposes no predicate for it.
const gsmMsgTypeLowerBound = message.MsgType(0x80)

// TestMsgTypeCoverage fails when github.com/free5gc/nas gains a 5GMM message
// type that msgTypeToMetricName does not map, which would otherwise silently
// start reporting the new message as UNKNOWN_GMM_MESSAGE.
func TestMsgTypeCoverage(t *testing.T) {
	for i := range int(gsmMsgTypeLowerBound) {
		msgType := message.MsgType(i)

		// An empty name means the upstream module does not know this type.
		if msgType.String() == "" {
			continue
		}

		if _, ok := msgTypeToMetricName[msgType]; !ok {
			t.Errorf("5GMM message type %s (%d) has no entry in msgTypeToMetricName",
				msgType, i)
		}
	}
}

// TestMsgTypeToMetricNameIsInjective guards against two message types sharing a
// metric name, which would silently merge their counters.
func TestMsgTypeToMetricNameIsInjective(t *testing.T) {
	seen := make(map[string]message.MsgType, len(msgTypeToMetricName))

	for msgType, name := range msgTypeToMetricName {
		if previous, ok := seen[name]; ok {
			t.Errorf("metric name %q is used by both %s and %s", name, previous, msgType)
			continue
		}
		seen[name] = msgType
	}
}

// TestMetricNameOf pins the label values that dashboards and alerting rules
// depend on. They must not follow message.MsgType.String() of the upstream
// module: those names are shorter and would break every existing query.
func TestMetricNameOf(t *testing.T) {
	testCases := []struct {
		msgType  message.MsgType
		expected string
	}{
		{message.MsgTypeAuthReq, AUTHENTICATION_REQUEST},
		{message.MsgTypeRegReq, REGISTRATION_REQUEST},
		{message.MsgTypeRegRej, REGISTRATION_REJECT},
		{message.MsgTypeDeregReqUETerm, DEREGISTRATION_REQUEST_UE_TERMINATED_DEREGISTRATION},
		{message.MsgTypeSvcRej, SERVICE_REJECT},
		{message.MsgTypeSecModeCmd, SECURITY_MODE_COMMAND},
		{message.MsgTypeCfgUpdateCmd, CONFIGURATION_UPDATE_COMMAND},
		{message.MsgTypeIdReq, IDENTITY_REQUEST},
		{message.MsgTypeNotif, NOTIFICATION},
		{message.MsgTypeStatus5GMM, STATUS_5GMM},
		{message.MsgTypeULNASTransport, UL_NAS_TRANSPORT},
		// 5GSM messages are not counted by this collector.
		{message.MsgTypePDUSessEstReq, UNKNOWN_GMM_MESSAGE},
		// Unassigned message type.
		{message.MsgType(0), UNKNOWN_GMM_MESSAGE},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expected, func(t *testing.T) {
			if name := metricNameOf(testCase.msgType); name != testCase.expected {
				t.Errorf("expected name: %q, output name: %q", testCase.expected, name)
			}
		})
	}
}

// TestGetMessageStrFromGmmMessage covers the message name and Cause5GMM
// extraction, including the nil shapes that reach this package from a deferred
// call whose message was never decoded.
func TestGetMessageStrFromGmmMessage(t *testing.T) {
	causeValue := uint8(3) // Cause5GMM_IllegalUE

	regRejWithCause := &message.RegRej{}
	regRejWithCause.Cause5GMM = &nasie.Cause5GMM{Value: causeValue}

	testCases := []struct {
		name          string
		msg           message.Message
		expectedName  string
		expectedCause *uint8
	}{
		{
			name:         "nil interface",
			msg:          nil,
			expectedName: UNKNOWN_GMM_MESSAGE,
		},
		{
			name:         "typed nil pointer",
			msg:          (*message.RegRej)(nil),
			expectedName: REGISTRATION_REJECT,
		},
		{
			name:         "message without cause IE",
			msg:          &message.RegReq{},
			expectedName: REGISTRATION_REQUEST,
		},
		{
			name:         "message with absent cause IE",
			msg:          &message.RegRej{},
			expectedName: REGISTRATION_REJECT,
		},
		{
			name:          "message with cause IE",
			msg:           regRejWithCause,
			expectedName:  REGISTRATION_REJECT,
			expectedCause: &causeValue,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ie := getMessageStrFromGmmMessage(testCase.msg)

			if ie.nasMessageType != testCase.expectedName {
				t.Errorf("expected name: %q, output name: %q", testCase.expectedName, ie.nasMessageType)
			}

			switch {
			case testCase.expectedCause == nil && ie.cause != nil:
				t.Errorf("expected no cause, output cause: %d", ie.cause.Value)
			case testCase.expectedCause != nil && ie.cause == nil:
				t.Errorf("expected cause: %d, output no cause", *testCase.expectedCause)
			case testCase.expectedCause != nil && ie.cause.Value != *testCase.expectedCause:
				t.Errorf("expected cause: %d, output cause: %d", *testCase.expectedCause, ie.cause.Value)
			}
		})
	}
}
