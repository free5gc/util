package nas

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
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

func IncrMetricsRcvNasMsg(msg *nas.Message, isStatusSuccess *bool, cause *string) {
	if IsNasMetricsEnabled() {
		nasMessageIe := getMessageStrFromGmmMessage(msg)
		metricCause := removeDigitSuffix(nasMessage.Cause5GMMToString(nasMessageIe.cause.Octet))
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
			errCause = removeDigitSuffix(nasMessage.Cause5GMMToString(cause5GMM))
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

type IeFromGmmMessage struct {
	nasMessageType string
	cause          nasType.Cause5GMM
}

func getMessageStrFromGmmMessage(msg *nas.Message) IeFromGmmMessage {
	ie := IeFromGmmMessage{nasMessageType: "Unknown gmm message"}

	if msg == nil || msg.GmmMessage == nil {
		return ie
	}

	switch {
	case msg.AuthenticationRequest != nil:
		ie.nasMessageType = AUTHENTICATION_REQUEST
	case msg.AuthenticationResponse != nil:
		ie.nasMessageType = AUTHENTICATION_RESPONSE
	case msg.AuthenticationResult != nil:
		ie.nasMessageType = AUTHENTICATION_RESULT
	case msg.AuthenticationFailure != nil:
		ie.nasMessageType = AUTHENTICATION_FAILURE
		ie.cause = msg.AuthenticationFailure.Cause5GMM
	case msg.AuthenticationReject != nil:
		ie.nasMessageType = AUTHENTICATION_REJECT
	case msg.RegistrationRequest != nil:
		ie.nasMessageType = REGISTRATION_REQUEST
	case msg.RegistrationAccept != nil:
		ie.nasMessageType = REGISTRATION_ACCEPT
	case msg.RegistrationComplete != nil:
		ie.nasMessageType = REGISTRATION_COMPLETE
	case msg.RegistrationReject != nil:
		ie.nasMessageType = REGISTRATION_REJECT
		ie.cause = msg.RegistrationReject.Cause5GMM
	case msg.ULNASTransport != nil:
		ie.nasMessageType = UL_NAS_TRANSPORT
	case msg.DLNASTransport != nil:
		ie.nasMessageType = DL_NAS_TRANSPORT
		ie.cause = *msg.DLNASTransport.Cause5GMM
	case msg.DeregistrationRequestUEOriginatingDeregistration != nil:
		ie.nasMessageType = DEREGISTRATION_REQUEST_UE_ORIGINATING_DEREGISTRATION
	case msg.DeregistrationAcceptUEOriginatingDeregistration != nil:
		ie.nasMessageType = DEREGISTRATION_ACCEPT_UE_ORIGINATING_DEREGISTRATION
	case msg.DeregistrationRequestUETerminatedDeregistration != nil:
		ie.nasMessageType = DEREGISTRATION_REQUEST_UE_TERMINATED_DEREGISTRATION
		ie.cause = *msg.DeregistrationRequestUETerminatedDeregistration.Cause5GMM
	case msg.DeregistrationAcceptUETerminatedDeregistration != nil:
		ie.nasMessageType = DEREGISTRATION_ACCEPT_UE_TERMINATED_DEREGISTRATION
	case msg.ServiceRequest != nil:
		ie.nasMessageType = SERVICE_REQUEST
	case msg.ServiceAccept != nil:
		ie.nasMessageType = SERVICE_ACCEPT
	case msg.ServiceReject != nil:
		ie.nasMessageType = SERVICE_REJECT
		ie.cause = msg.ServiceReject.Cause5GMM
	case msg.ConfigurationUpdateCommand != nil:
		ie.nasMessageType = CONFIGURATION_UPDATE_COMMAND
	case msg.ConfigurationUpdateComplete != nil:
		ie.nasMessageType = CONFIGURATION_UPDATE_COMPLETE
	case msg.IdentityRequest != nil:
		ie.nasMessageType = IDENTITY_REQUEST
	case msg.IdentityResponse != nil:
		ie.nasMessageType = IDENTITY_RESPONSE
	case msg.Notification != nil:
		ie.nasMessageType = NOTIFICATION
	case msg.NotificationResponse != nil:
		ie.nasMessageType = NOTIFICATION_RESPONSE
	case msg.SecurityModeCommand != nil:
		ie.nasMessageType = SECURITY_MODE_COMMAND
	case msg.SecurityModeComplete != nil:
		ie.nasMessageType = SECURITY_MODE_COMPLETE
	case msg.SecurityModeReject != nil:
		ie.nasMessageType = SECURITY_MODE_REJECT
		ie.cause = msg.SecurityModeReject.Cause5GMM
	case msg.SecurityProtected5GSNASMessage != nil:
		ie.nasMessageType = SECURITY_PROTECTED_5GS_NAS_MESSAGE
	case msg.Status5GMM != nil:
		ie.nasMessageType = STATUS_5GMM
		ie.cause = msg.Status5GMM.Cause5GMM
	}
	return ie
}
