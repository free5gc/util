// Package nas: NAS message metrics, built against the new
// github.com/free5gc/nas API (ie/message packages).
package nas

import (
	"reflect"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"

	nasie "github.com/free5gc/nas/ie"
	"github.com/free5gc/nas/message"
	"github.com/free5gc/util/metrics/utils"
)

var suffixRe = regexp.MustCompile(`\s*\(\d+\)$`)

var (
	NasMsgRcvCounter  *prometheus.CounterVec
	NasMsgSentCounter *prometheus.CounterVec
)

func GetNasHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	NasMsgRcvCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      NAS_MSG_RCV_COUNTER_NAME,
			Help:      NAS_MSG_RCV_COUNTER_DESC,
		},
		[]string{NAME_LABEL, STATUS_LABEL, CAUSE_LABEL},
	)

	collectors = append(collectors, NasMsgRcvCounter)

	NasMsgSentCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      NAS_MSG_SENT_COUNTER_NAME,
			Help:      NAS_MSG_SENT_COUNTER_DESC,
		},
		[]string{NAME_LABEL, STATUS_LABEL, CAUSE_LABEL},
	)

	collectors = append(collectors, NasMsgSentCounter)

	return collectors
}

func removeDigitSuffix(s string) string {
	return suffixRe.ReplaceAllString(s, "")
}

func IncrMetricsRcvNasMsg(msg message.Message, isStatusSuccess *bool, cause *string) {
	if IsNasMetricsEnabled() {
		nasMessageIe := getMessageStrFromGmmMessage(msg)
		metricCause := ""
		if nasMessageIe.cause != nil {
			metricCause = removeDigitSuffix(nasMessageIe.cause.String())
		}
		metricStatus := utils.FailureMetric

		if cause != nil && *cause != "" {
			metricCause = *cause
		}

		if isStatusSuccess != nil && *isStatusSuccess {
			metricStatus = utils.SuccessMetric
		}

		NasMsgRcvCounter.With(prometheus.Labels{
			NAME_LABEL:   nasMessageIe.nasMessageType,
			STATUS_LABEL: metricStatus,
			CAUSE_LABEL:  metricCause,
		}).Inc()
	}
}

func IncrMetricsSentNasMsgs(msgType string, isStatusSuccess *bool, cause5GMM uint8, otherCause *string) {
	if IsNasMetricsEnabled() {
		errCause := ""

		if cause5GMM != 0 {
			gmmCause := nasie.Cause5GMM{Value: cause5GMM}
			errCause = removeDigitSuffix(gmmCause.String())
		} else if otherCause != nil {
			errCause = *otherCause
		}

		metricStatus := utils.FailureMetric

		if isStatusSuccess != nil && *isStatusSuccess {
			metricStatus = utils.SuccessMetric
		}

		NasMsgSentCounter.With(prometheus.Labels{
			NAME_LABEL:   msgType,
			STATUS_LABEL: metricStatus,
			CAUSE_LABEL:  errCause,
		}).Inc()
	}
}

// msgTypeToMetricName pins the "name" label of the NAS metrics to the constants
// in types.go instead of letting it follow message.MsgType.String(), whose
// values belong to github.com/free5gc/nas and may change without notice here.
// Every 5GMM message type of that module must have an entry; msgTypeCoverage in
// message_test.go fails when upstream adds one that is missing below.
var msgTypeToMetricName = map[message.MsgType]string{
	message.MsgTypeRegReq:                      REGISTRATION_REQUEST,
	message.MsgTypeRegAccept:                   REGISTRATION_ACCEPT,
	message.MsgTypeRegComplete:                 REGISTRATION_COMPLETE,
	message.MsgTypeRegRej:                      REGISTRATION_REJECT,
	message.MsgTypeDeregReqUEOrig:              DEREGISTRATION_REQUEST_UE_ORIGINATING_DEREGISTRATION,
	message.MsgTypeDeregAcceptUEOrig:           DEREGISTRATION_ACCEPT_UE_ORIGINATING_DEREGISTRATION,
	message.MsgTypeDeregReqUETerm:              DEREGISTRATION_REQUEST_UE_TERMINATED_DEREGISTRATION,
	message.MsgTypeDeregAcceptUETerm:           DEREGISTRATION_ACCEPT_UE_TERMINATED_DEREGISTRATION,
	message.MsgTypeSvcReq:                      SERVICE_REQUEST,
	message.MsgTypeSvcRej:                      SERVICE_REJECT,
	message.MsgTypeSvcAccept:                   SERVICE_ACCEPT,
	message.MsgTypeCtrlPlaneSvcReq:             CONTROL_PLANE_SERVICE_REQUEST,
	message.MsgTypeNwSliceSpecificAuthCmd:      NETWORK_SLICE_SPECIFIC_AUTHENTICATION_COMMAND,
	message.MsgTypeNwSliceSpecificAuthComplete: NETWORK_SLICE_SPECIFIC_AUTHENTICATION_COMPLETE,
	message.MsgTypeNwSliceSpecificAuthResult:   NETWORK_SLICE_SPECIFIC_AUTHENTICATION_RESULT,
	message.MsgTypeCfgUpdateCmd:                CONFIGURATION_UPDATE_COMMAND,
	message.MsgTypeCfgUpdateComplete:           CONFIGURATION_UPDATE_COMPLETE,
	message.MsgTypeAuthReq:                     AUTHENTICATION_REQUEST,
	message.MsgTypeAuthRsp:                     AUTHENTICATION_RESPONSE,
	message.MsgTypeAuthRej:                     AUTHENTICATION_REJECT,
	message.MsgTypeAuthFailure:                 AUTHENTICATION_FAILURE,
	message.MsgTypeAuthResult:                  AUTHENTICATION_RESULT,
	message.MsgTypeIdReq:                       IDENTITY_REQUEST,
	message.MsgTypeIdRsp:                       IDENTITY_RESPONSE,
	message.MsgTypeSecModeCmd:                  SECURITY_MODE_COMMAND,
	message.MsgTypeSecModeComplete:             SECURITY_MODE_COMPLETE,
	message.MsgTypeSecModeRej:                  SECURITY_MODE_REJECT,
	message.MsgTypeStatus5GMM:                  STATUS_5GMM,
	message.MsgTypeNotif:                       NOTIFICATION,
	message.MsgTypeNotifRsp:                    NOTIFICATION_RESPONSE,
	message.MsgTypeULNASTransport:              UL_NAS_TRANSPORT,
	message.MsgTypeDLNASTransport:              DL_NAS_TRANSPORT,
	message.MsgTypeRelayKeyReq:                 RELAY_KEY_REQUEST,
	message.MsgTypeRelayKeyAccept:              RELAY_KEY_ACCEPT,
	message.MsgTypeRelayKeyRej:                 RELAY_KEY_REJECT,
	message.MsgTypeRelayAuthReq:                RELAY_AUTHENTICATION_REQUEST,
	message.MsgTypeRelayAuthRsp:                RELAY_AUTHENTICATION_RESPONSE,
}

// metricNameOf returns the "name" label value for a message type, falling back
// to UNKNOWN_GMM_MESSAGE for GSM messages and for types this package does not
// know about.
func metricNameOf(msgType message.MsgType) string {
	if name, ok := msgTypeToMetricName[msgType]; ok {
		return name
	}

	return UNKNOWN_GMM_MESSAGE
}

type IeFromGmmMessage struct {
	nasMessageType string
	cause          *nasie.Cause5GMM
}

func getMessageStrFromGmmMessage(msg message.Message) IeFromGmmMessage {
	ie := IeFromGmmMessage{nasMessageType: UNKNOWN_GMM_MESSAGE}

	if msg == nil {
		return ie
	}

	ie.nasMessageType = metricNameOf(msg.MsgType())

	// A typed-nil pointer stored in the interface passes the msg == nil check
	// above and survives MsgType(), but would panic on the field accesses
	// below. Counting it under its message name with no cause is enough.
	if v := reflect.ValueOf(msg); v.Kind() == reflect.Pointer && v.IsNil() {
		return ie
	}

	// Only the 5GMM messages carrying a Cause5GMM IE need a case here, the
	// message name is already resolved above.
	switch m := msg.(type) {
	case *message.AuthFailure:
		ie.cause = m.Cause5GMM
	case *message.RegRej:
		ie.cause = m.Cause5GMM
	case *message.DLNASTransport:
		ie.cause = m.Cause5GMM
	case *message.DeregReqUETerm:
		ie.cause = m.Cause5GMM
	case *message.SvcRej:
		ie.cause = m.Cause5GMM
	case *message.SecModeRej:
		ie.cause = m.Cause5GMM
	case *message.Status5GMM:
		ie.cause = m.Cause5GMM
	}
	return ie
}
